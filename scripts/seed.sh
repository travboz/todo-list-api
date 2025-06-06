#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <BearerToken>"
  exit 1
fi

TOKEN="$1"
ENDPOINT="http://localhost:8000/api/v1/tasks/create"

titles=(
  "Send welcome email"
  "Clean bedroom"
  "Daily team meeting"
  "Plan travel itinerary"
  "Review budget plan"
  "Write Go tutorial"
  "Book flights for travel"
  "Organize email inbox"
  "Prepare meeting agenda"
  "Clean kitchen"
  "Submit project report"
  "Pay monthly budget"
  "Fix Go code bugs"
  "Send project update email"
  "Travel insurance renewal"
  "Team meeting follow-up"
  "Deep clean bathroom"
  "Read Go concurrency chapter"
  "Email client proposal"
  "Plan office travel"
  "Budget review with finance"
  "Refactor Go service"
  "Schedule meeting with team"
  "Clean living room"
  "Write email newsletter"
  "Update travel documents"
  "Go performance tuning"
  "Meeting notes compilation"
  "Pay credit card bill"
  "Send birthday email"
  "Plan marketing travel"
  "Clean garage"
  "Write Go unit tests"
  "Review email templates"
  "Fix meeting room projector"
  "Update budget spreadsheet"
  "Organize travel bags"
  "Go module version bump"
  "Team building meeting"
  "Clean dining room"
  "Email follow-up to client"
  "Prepare travel checklist"
  "Write Go benchmarks"
  "Budget allocation meeting"
  "Clean office desks"
  "Send holiday email"
  "Book travel accommodations"
  "Review Go code standards"
  "Plan meeting refreshments"
  "Pay rent"
  "Email monthly report"
  "Travel expense submission"
)

descriptions=(
  "Draft and send the welcome email to new subscribers."
  "Vacuum the floor and wipe surfaces in the bedroom."
  "Join the daily standup meeting at 9 AM."
  "Create a travel plan including destinations and accommodation."
  "Review and adjust the monthly budget plan."
  "Write a tutorial on Go slices and maps."
  "Book flights and hotels for upcoming travel."
  "Organize and archive old emails in inbox."
  "Prepare the agenda for next week's team meeting."
  "Clean countertops, mop floor, and empty trash in kitchen."
  "Complete and submit the final project report."
  "Pay electricity, water, and internet bills this month."
  "Fix bugs found in the Go backend service."
  "Send an update email about the current project status."
  "Renew travel insurance before the trip."
  "Send follow-up emails after the team meeting."
  "Deep clean tiles, grout, and fixtures in bathroom."
  "Read the chapter on concurrency in the Go programming book."
  "Email the client the updated proposal documents."
  "Plan travel arrangements for the office retreat."
  "Meet with finance to review the annual budget."
  "Refactor the Go microservice for better performance."
  "Schedule a meeting with the development team."
  "Clean and tidy up the living room."
  "Write and send the monthly email newsletter."
  "Update passports and travel visas before trip."
  "Optimize Go application performance and memory."
  "Compile notes and action items from meetings."
  "Pay credit card bill before due date."
  "Send birthday greeting emails to clients."
  "Plan travel for upcoming marketing campaign."
  "Organize and clean out the garage."
  "Write unit tests for Go codebase."
  "Review and improve email marketing templates."
  "Fix the projector in the main meeting room."
  "Update budget spreadsheets with new data."
  "Pack and organize bags for travel."
  "Bump Go module versions to latest releases."
  "Arrange team building meeting and activities."
  "Clean the dining room and polish furniture."
  "Send follow-up email to potential client."
  "Prepare a checklist for the travel trip."
  "Write benchmark tests for Go packages."
  "Discuss budget allocation in next meeting."
  "Clean office desks and dispose of trash."
  "Send holiday greetings emails."
  "Book hotels and transport for travel."
  "Review and enforce Go coding standards."
  "Arrange refreshments for next meeting."
  "Pay monthly rent on time."
  "Send monthly progress report via email."
  "Submit travel expense claims and receipts."
)

for i in {0..49}; do
  curl -s -X POST "$ENDPOINT" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
      \"title\": \"${titles[$i]}\",
      \"description\": \"${descriptions[$i]}\"
    }"

  echo ""
done
