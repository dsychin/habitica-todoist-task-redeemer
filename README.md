# Habitica Todoist Task Redeemer

[![Go](https://github.com/dsychin/habitica-todoist-task-redeemer/actions/workflows/go.yml/badge.svg)](https://github.com/dsychin/habitica-todoist-task-redeemer/actions/workflows/go.yml)

Integration to add and mark the task complete to Habitica when a task is completed in Todoist.

This can be used so that your todos can be managed in Todoist without losing out on experience, etc in Habitica.

## Features

1. Automatically add a task to Habitica with Easy difficulty and mark it as complete when a task has been completed in Todoist.

## Prerequisites

- A GitHub account

## Instructions

### First time setup

[![](https://www.netlify.com/img/deploy/button.svg)](https://app.netlify.com/start/deploy?repository=https://github.com/dsychin/habitica-todoist-task-redeemer)

1. Use the deploy to Netlify button above and follow the instructions to login with your GitHub account.
2. Fill in your Habitica API token and user id when prompted.
3. After deploying, take note of your site URL in the Site Overview page. It should look something like `somerandomvalue.netlify.app`.
4. Go to the [Todoist App Management Console](https://developer.todoist.com/appconsole.html) and create a new app.
5. Give it an app name and leave the service url blank.
6. Under the "Test token" section, press "Create test token". This will make it so that events from your Todoist account will send a webhook.
7. Under the "Webhook" section, set the Webhook Callback URL to your Netlify site URL with the path `/.netlify/functions/redeemer`. For example, if your Netlify URL is `somerandomvalue.netlify.app`, then the URL you should put should be `https://somerandomvalue.netlify.app/.netlify/functions/redeemer`.
8. For watched events, check `item:completed` and `item:uncompleted`.
9. Save your Todoist app configuration.

### Usage

1. In Todoist, mark an item as completed.
2. In Habitica, you should see an item has been added in your ToDos as completed. You may need to change your filter to show completed ToDos.

### Debugging issues

1. In Netlify's site dashboard, go to the Functions tab and click on the `redeemer` function.
2. This will show the logs for that function.
3. Mark an item as completed in Todoist and view the logs in Netlify.

## Limitations

- Currently does not verify Todoist's webhook signature so there is a risk of someone sending similar data to your endpoint.
- Accidentally marking something and unmarking it does not undo the task in Habitica. It will need to be unchecked manually.
- Task difficulty is set to Easy and cannot be changed unless you change it in the code directly.
- Currently does not check who completes the task or who is assigned in Todoist, so if you are using this in a shared project it might add all tasks to your Habitica account.

## Similar projects

- [Todoist Sync](https://habitica.fandom.com/wiki/Todoist_Sync)
- [Habitica Todo](https://habitica.fandom.com/wiki/Habitica-Todo)
