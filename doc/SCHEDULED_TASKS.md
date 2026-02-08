# Scheduled Tasks in JGBot

JGBot supports two types of scheduled tasks: **Cron Jobs** and **Timers**. These allow the bot to perform actions or send messages autonomously at specific times or intervals.

## Overview

Scheduled tasks work by sending a trigger message to the bot's agent. The agent then processes this message and responds according to its instructions, just like it would to a user message.

### Core Tools

The following tools are available to the agent (if enabled in the session configuration):

- `cron`: For recurring tasks using cron expressions.
- `timer`: For one-time tasks (back-counters) or specific alarms.

## Cron Jobs

Cron jobs are used for tasks that repeat according to a schedule.

### Actions
- `list`: Shows all active cron jobs for the current session.
- `read`: Shows detailed information about a specific cron job.
- `add`: Creates a new recurring task.
- `remove`: Deletes an existing cron job.

### Cron Expressions
JGBot uses a standard 5 field cron expression format:
`Minutes | Hours | Day of Month | Month | Day of Week`

**Examples:**
- `0 12 * * ?`: Every day at 12:00 PM.
- `*/5 * * * *`: Every 5 minutes.
- `0 9 * * 1`: Every Monday at 9:00 AM.

## Timers & Alarms

Timers are used for one-time events.

### Types
- **Timeout**: Triggers after a certain duration (e.g., "in 30 minutes").
- **Alarm**: Triggers at a specific time (e.g., "at 8:00 AM").

### Actions
- `list`: Shows all pending timers for the current session.
- `read`: Shows details of a pending timer.
- `add`: Sets a new timeout or alarm.
- `remove`: Cancels a pending timer.

## How it Works

When a task triggers:
1. The system generates a "virtual" message containing the task's Name, Description, and Message.
2. This message is injected into the session's conversation history.
3. The AI Agent is activated and receives this message.
4. The Agent responds to the message (e.g., by sending a reminder, performing a search, or executing a skill).

### Example Workflow
1. **User**: "Remind me to water the plants every day at 8 AM."
2. **Bot** (via `cron.add`): Sets a job named `PlantReminder` with schedule `0 8 * * *`.
3. **Execution** (at 8 AM):
   - System: `CRON EXECUTION: PlantReminder ... MESSAGE: Remind the user to water the plants.`
   - Agent: "Hey! Don't forget to water your plants today. 🌱"

## Persistence

All scheduled tasks are persisted in the `config/` directory:
- `config/cron.json`
- `config/timers.json`

This ensures that tasks are not lost if the bot is restarted.

## Message Reactions

While not a scheduling feature, the `message_reaction` tool is another recent addition that allows the agent to interact visually with messages.

### Usage
- `message_id`: The ID of the message to react to.
- `reaction`: The emoji string (e.g., "👍", "❤️").
