---
name: write-docs
description: Write or edit OpenGovMail documentation (README.md, CONTRIBUTING.md, chart READMEs and pages under docs/) following the project's style guide, topic types and preferred terms. Use when creating or changing any Markdown documentation in this repository.
---

# Write documentation

[docs/style-guide.md](../../../docs/style-guide.md) is the source of truth. It builds on the [Google developer documentation style guide](https://developers.google.com/style), and where they disagree, the project's guide wins. Read it before you write.

## Steps

1. **Pick the topic type first.** Each page is exactly one type, and it goes in a fixed place:

   | The page... | Type | Where |
   |---|---|---|
   | explains an idea | Concept | `docs/concepts/` |
   | gives steps for one task | How-to | `docs/how-to/`, or `docs/user-guide/` for end users |
   | lists facts to look up | Reference | `docs/reference/` |
   | lists problems and fixes | Troubleshooting | `docs/troubleshooting.md` |
   | defines a term | Glossary | `docs/glossary.md` |

   If the content mixes types, split it into separate pages and link them. Each type follows its template from [The Good Docs Project](https://gitlab.com/tgdp/templates). Use a page of the same type in `docs/` as a model.

2. **Follow the writing rules:**
   - Second person ("you"), present tense, active voice.
   - Sentence case for headings.
   - Lead with what the reader does, then give the reason in one sentence.
   - One task per how-to, with no more than about 10 steps.
   - Commands in fenced code blocks that can be copied and run as they are, with the parts to replace in `<angle brackets>`.
   - Use bold sparingly, and don't use em dashes.
   - When the page states a version, say which chart release it's for.
   - Diagrams are Mermaid blocks. Give every image alt text, and don't rely on colour alone.

3. **Use the preferred terms** from the glossary table in the style guide. For example, write "ThunderID" (not "Thunder"), "mail app" (not "MUA"), "organization unit" (not "OU" in running text), `example.com` for domains and `12.34.56.78` for IP addresses. Define any new term in [docs/glossary.md](../../../docs/glossary.md).

4. **Document only Helm on Kubernetes.** It's the only supported install method.

5. **Link the page in:**
   - Add a new page to [docs/README.md](../../../docs/README.md).
   - Make link text say where the link goes, never "click here".
   - Use relative links between files in the repository.

6. **Check the links.** Every relative link must point to a file that exists. A quick check from the repository root, for the page you changed:

   ```bash
   page=docs/<path>.md
   grep -oE '\]\([^)#]+' "$page" | sed 's/](//' | grep -vE '^(https?|mailto):' | while read -r link; do
     [ -e "$(dirname "$page")/$link" ] || echo "broken: $link"
   done
   ```

7. **Update the docs with the change they describe,** in the same pull request. If you change the style guide itself, add a row to its decision log.

## Done when

- [ ] The page is one topic type, in the right folder, and follows that type's template
- [ ] It follows the writing rules and uses the preferred terms
- [ ] New pages are listed in `docs/README.md`, and all relative links resolve
- [ ] No real domains, IP addresses, credentials or personal data appear
