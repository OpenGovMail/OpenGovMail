import json, os, sys, time, urllib.error, urllib.parse, urllib.request
from base64 import b64encode

# Registers one OAuth client in ThunderID, or updates it in place, authenticating as
# raven's m2m client. The chart passes the application body in APP_JSON; this fills
# in what is only known at runtime: the OU and flow ids, and the client secret.

BASE = os.environ["THUNDER_URL"]

def call(path, token=None, data=None, method=None, basic=None):
    req = urllib.request.Request(BASE + path, method=method)
    if data is not None:
        req.data = data if isinstance(data, bytes) else json.dumps(data).encode()
        req.add_header("Content-Type", "application/json")
    if token:
        req.add_header("Authorization", "Bearer " + token)
    if basic:
        req.add_header("Authorization", "Basic " + b64encode(basic.encode()).decode())
    with urllib.request.urlopen(req, timeout=30) as r:
        body = r.read()
        return r.status, (json.loads(body) if body else {})

# ThunderID can still be starting when this hook runs.
for _ in range(60):
    try:
        call("/.well-known/openid-configuration")
        break
    except Exception:
        time.sleep(5)
else:
    sys.exit("ThunderID never became reachable")

form = urllib.parse.urlencode({
    "grant_type": "client_credentials",
    "resource": os.environ["RESOURCE"],
    "scope": "system",
}).encode()
req = urllib.request.Request(BASE + "/oauth2/token", data=form, method="POST")
req.add_header("Authorization", "Basic " + b64encode(
    f"{os.environ['RAVEN_CLIENT_ID']}:{os.environ['RAVEN_CLIENT_SECRET']}".encode()).decode())
with urllib.request.urlopen(req, timeout=30) as r:
    token = json.loads(r.read())["access_token"]

# Lists default to 30 items, so read every page.
def list_all(path, key):
    items = []
    while True:
        _, page = call(f"{path}?limit=100&offset={len(items)}", token)
        batch = page.get(key) or []
        items += batch
        if not batch or len(items) >= page.get("totalResults", 0):
            return items

# Resolve ids from handles. The UUIDs are install data, the handles are not.
# Neither endpoint filters server side, so both are matched here.
ou_handle = os.environ["OU_HANDLE"]
ou = next((o for o in list_all("/organization-units", "organizationUnits") if o.get("handle") == ou_handle), None)
if not ou:
    sys.exit(f"no organization unit with handle {ou_handle}")

want = os.environ["AUTH_FLOW_HANDLE"]
flow_id = None
for f in list_all("/flows", "flows"):
    if f.get("handle") != want:
        continue
    # flowType is only on the individual resource, not the list.
    _, full = call("/flows/" + f["id"], token)
    if full.get("flowType") == "AUTHENTICATION":
        flow_id = f["id"]
        break
if not flow_id:
    sys.exit(f"no AUTHENTICATION flow with handle {want}")

client_id = os.environ["CLIENT_ID"]
existing = next((a for a in list_all("/applications", "applications")
                 if a.get("clientId") == client_id), None)

app = json.loads(os.environ["APP_JSON"])
app["ouId"] = ou["id"]
app["authFlowId"] = flow_id
oauth = app["inboundAuthConfig"][0]["config"]

# ThunderID rejects an attribute that none of the allowed user types declares, and
# user types are defined per install. With PRUNE_USER_ATTRIBUTES set, drop the ones
# they lack, but fail if any listed there as required is missing.
if os.environ.get("PRUNE_USER_ATTRIBUTES"):
    by_name = {t["name"]: t["id"] for t in list_all("/user-types", "types")}
    declared = set()
    for name in app.get("allowedUserTypes") or []:
        if name not in by_name:
            sys.exit(f"no user type named {name}")
        _, user_type = call("/user-types/" + by_name[name], token)
        declared |= set((user_type.get("schema") or {}).keys())
    missing = [a for a in os.environ["PRUNE_USER_ATTRIBUTES"].split(",") if a not in declared]
    if missing:
        sys.exit(f"user types {app['allowedUserTypes']} declare no {', '.join(missing)} attribute")
    keep = lambda attrs: [a for a in attrs if a in declared]
    app["userAttributes"] = keep(app.get("userAttributes") or [])
    for holder in (oauth.get("token", {}).get("idToken"), oauth.get("token", {}).get("accessToken"),
                   oauth.get("userInfo")):
        if holder and "userAttributes" in holder:
            holder["userAttributes"] = keep(holder["userAttributes"])
    for scope, claims in (oauth.get("scopeClaims") or {}).items():
        oauth["scopeClaims"][scope] = keep(claims)
    print("user attributes:", ", ".join(app["userAttributes"]))
oauth["clientId"] = client_id
oauth["clientSecret"] = os.environ["CLIENT_SECRET"]

# PUT takes the same body as create, so the redirect URI and secret
# follow the chart's values on every upgrade.
try:
    if existing:
        status, _ = call("/applications/" + existing["id"], token, data=app, method="PUT")
        print("updated", client_id, status)
    else:
        status, _ = call("/applications", token, data=app, method="POST")
        print("registered", client_id, status)
except urllib.error.HTTPError as e:
    # A clash means a concurrent run created it; the retry updates it.
    sys.exit(f"{e.code}: {e.read().decode()}")
