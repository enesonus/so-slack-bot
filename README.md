# SO Slack Bot

`A Slack bot for reporting Stack Overflow questions`

SO Slack is a Slack bot that helps software projects improve developer experience by sending Slack messages about Stack Overflow questions.

Keeping an eye on Stack Overflow questions helps makes projects more reachable for developers. Using Stack Overflow is a quick way to spot and fix common problems, get real feedback, and make your software easier to use. Having lower response time on Stack Overflow results in happier developers and an increasingly better and bigger project.

## Usage

If you want to use the bot at your Slack Workspace click:

<a href="https://slack.com/oauth/v2/authorize?client_id=5970139268528.5949679575364&scope=channels:history,chat:write,groups:history,im:history,mpim:history,team:read,users:read,app_mentions:read,channels:join,channels:read&user_scope="><img alt="Add to Slack" height="40" width="139" src="https://platform.slack-edge.com/img/add_to_slack.png" srcSet="https://platform.slack-edge.com/img/add_to_slack.png 1x, https://platform.slack-edge.com/img/add_to_slack@2x.png 2x" /></a>

If you want to host the bot on your own, navigate to [Deployment](#deployment)

### Commands to use SO Slack bot

**Reminder**: Make sure to add the bot in the channel where you plan to use it.

1. `soslack_set_so_channel`: Sets a channel for the bot to post messages in.
2. `soslack_remove_so_channel`: Stops the bot from sending messages to a channel.
3. `soslack_add_tag <tag>`: Links a tag with a channel for targeted message notifications.
4. `soslack_show_tags`: Displays all tags linked to the current channel for notifications.

## Deployment

```bash
Work in Progress
```
