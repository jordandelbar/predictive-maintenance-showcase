HEALTH_CHECK_URL="http://localhost:4000/health"
ROUTE="http://localhost:4000/v1/threshold"

while true; do
  echo "Checking health at: $HEALTH_CHECK_URL"

  if curl -X GET "$HEALTH_CHECK_URL"; then
    if curl -i -X POST "$ROUTE" -d '{"machine_id": 7, "threshold": 0.02}'; then
      break
    else
      echo "Failed to send POST request. Retrying in 3 seconds..."
      sleep 3
    fi
  else
    echo "Service is not healthy, retrying in 3 seconds..."
    sleep 3
  fi
done
