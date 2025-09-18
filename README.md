# Repo Recap

Repo Recap gives people a daily or weekly digest of recent changes in GitHub repositories of interest.

The digest can be delivered via Slack, Discord, or email.

The digests will be generated from an LLM.

## Requirements

- Have a YAML config for adjusting the app
- Have a Dockerfile and publish images in this repo's registry
- Can choose all the GitHub repositories of interest - have separate messages for separate repositories
- Can choose between daily or weekly updates - if there's nothing new, no messages are sent
- Can pass in multiple Slack and Discord channels (emails will be supported later, if at all)
- Can choose between Gemini or ChatGPT for digest creation
- Can provide a custom message to the LLM regarding message formatting
