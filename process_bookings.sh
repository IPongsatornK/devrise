#!/usr/bin/env sh
set -e

INPUT="${1:-bookings.psv}"
OUTPUT="${2:-report.txt}"

if [ ! -f "$INPUT" ]; then
  echo "Input file not found: $INPUT" >&2
  exit 1
fi

now="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
generated="$(date -u +"%Y-%m-%d %H:%M:%S UTC")"

awk -F'|' -v now="$now" -v generated="$generated" '
NR==1 { next }
{
  status=$3
  amount=$4 + 0
  deadline=$6
  if (status == "CONFIRMED") {
    confirmed++
    revenue += amount
  } else if (status == "PENDING") {
    pending++
    if (deadline < now) {
      overdue[++count] = sprintf("  %s  %s  deadline: %s", $1, $2, deadline)
    }
  } else if (status == "CANCELLED") {
    cancelled++
  }
}
END {
  print "=== Booking Report (generated: " generated ") ===\n"
  print "Status Summary:"
  printf "  CONFIRMED : %d\n  PENDING   : %d\n  CANCELLED : %d\n", confirmed, pending, cancelled
  printf "\nTotal Revenue (CONFIRMED): %.2f THB\n\n", revenue
  print "Expired PENDING Bookings (payment overdue):"
  if (count == 0) {
    print "  None"
  } else {
    for (i = 1; i <= count; i++) {
      print overdue[i]
    }
  }
}
' "$INPUT" > "$OUTPUT"
