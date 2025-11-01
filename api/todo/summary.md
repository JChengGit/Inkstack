# ROLE
- You are a senior engineer and technical writer. 

# GOAL
- Produce a single Markdown file that contains:
  - A concise architecture overview of the function.
  - A machine-readable catalog of components/features for fast onboarding in future iterations.

# CONSTRAINTS
- Document only what is present in the provided code or context. If any information is missing, write: `Unknown/Needs verification`.
- Use standard Markdown formatting (titles, sections, bullet points). No advanced syntax, no icons.
- Ignore vendor/compiled folders (node_modules, .husky, .vscode, build, coverage, etc).

# REFERENCE
- `/local/PO-import-task.md` - initial description of the PO import refactoring task.
- `/local/bulkJobService.md` – a catalog explaining functions in `bulkJobService.ts`. Use this to understand the service without reading the entire large file.

# TASK
- Analyze the PO import function in this project.
- Write the result into `/local/PO-import-feature.md`.
