# When you run software engineering tasks, you MUST follow these principles
- Document files that are prefixed PLAN contain work plans for you. There could be an overarching PLAN.md document, and if there is that contains the highest level plan for your work. Documents that are named PLAN_{{FEATURE}}.md are plans that are focused on specific features or parts of the over all plan. If there is an overarching plan, it may control the dependencies between the different parts.
- If there is no overarching PLAN.md document, then you will work with just the specific PLAN_{{FEATURE}}.md files, and if there are several such files, you must think very hard about if there are dependencies between them that have consequences for what to do in what order.
- A README file, if it exists, will describe the overall purpose, usage, and technology aspects or dependencies of the project.
- If the PLAN documents are not loaded into the current context, follow this routine;
 1. Start with the overarching document if it exists (e.g. PLAN.md)
 2. Analyze the overarching plan, and conclude if there is a PLAN_{{FEATURE}}.md file that should be loaded in order to get the most relevant next part of the plan to work with.
 3. Load the next PLAN_{{FEATURE}}.md document
 4. Analyze and summarize the next step according to the overarching and specific part of the plan.


