#/bin/bash

GENERATED_REQ_RATE_HIGH="${GENERATED_REQ_RATE_HIGH:-15}"
GENERATED_REQ_RATE_HIGH_SECONDS="${GENERATED_REQ_RATE_HIGH_SECONDS-10}"
GENERATED_REQ_RATE_LOW="${GENERATED_REQ_RATE_LOW-1}"
GENERATED_REQ_RATE_LOW_SECONDS="${GENERATED_REQ_RATE_LOW_SECONDS-5}"
THRESHOLD_REQ_RATE_HIGH="${THRESHOLD_REQ_RATE_HIGH-13}"
THRESHOLD_REQ_RATE_HIGH_SECONDS="${THRESHOLD_REQ_RATE_HIGH_SECONDS-8}"
ANALYSIS_SECONDS="${ANALYSIS_SECONDS-15}"
CYCLES="${CYCLES-1}"
EXPECTED_HIGH_TRAFFIC_ALERTS="${EXPECTED_HIGH_TRAFFIC_ALERTS-1}"
EXPECTED_HIGH_TRAFFIC_RECOVERIES="${EXPECTED_HIGH_TRAFFIC_RECOVERIES-1}"

which -s timeout
if [[ $? -ne 0 ]] ; then
  echo "\"timeout\" command not found, please install it or fix your path"
  exit 1
fi

echo """=== alerting logic test ===

Generate synthetic logs:
1. Generate logs with request rate of $GENERATED_REQ_RATE_HIGH/s for $GENERATED_REQ_RATE_HIGH_SECONDS seconds, then
2. Generate logs with request rate of $GENERATED_REQ_RATE_LOW/s for $GENERATED_REQ_RATE_LOW_SECONDS seconds, then
3. Repeat for $CYCLES cycle/s;

Run the logmonitor for $ANALYSIS_SECONDS seconds.
Alert if request rate is > $THRESHOLD_REQ_RATE_HIGH/s for $THRESHOLD_REQ_RATE_HIGH_SECONDS seconds.

We expect to see \"High traffic generated an alert\" $EXPECTED_HIGH_TRAFFIC_ALERTS time/s.
We expect to see \"High traffic alert recovered\" $EXPECTED_HIGH_TRAFFIC_RECOVERIES time/s.
"""

echo "starting alerting logic test...\n"

parent_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

logfile=`mktemp`
analysislog=`mktemp`

$parent_path/../bin/loggenerator \
    -f $logfile \
    -R $GENERATED_REQ_RATE_HIGH \
    -r $GENERATED_REQ_RATE_LOW \
    -T $GENERATED_REQ_RATE_HIGH_SECONDS \
    -t $GENERATED_REQ_RATE_LOW_SECONDS \
    -c $CYCLES > /dev/null &

timeout --foreground $ANALYSIS_SECONDS $parent_path/../bin/logmonitor \
    -f $logfile \
    -r $THRESHOLD_REQ_RATE_HIGH \
    -t $THRESHOLD_REQ_RATE_HIGH_SECONDS > $analysislog

retcode=0

grep -c "High traffic generated an alert" $analysislog | grep -q "^$EXPECTED_HIGH_TRAFFIC_ALERTS$"
if [[ $? -ne 0 ]] ; then
 echo "\"High traffic generated an alert\" not found in analysis exactly $EXPECTED_HIGH_TRAFFIC_ALERTS time/s!"
 retcode=1
fi

grep -c "High traffic alert recovered" $analysislog | grep -q "^$EXPECTED_HIGH_TRAFFIC_RECOVERIES$"
if [[ $? -ne 0 ]] ; then
 echo "\"High traffic alert recovered\" not found in analysis exactly $EXPECTED_HIGH_TRAFFIC_RECOVERIES time/s!"
 retcode=1
fi

if [[ $retcode -eq 0 ]] ; then
  echo "test successfully passed!"
fi

echo
echo "logmonitor output was:"
cat $analysislog

rm $logfile
rm $analysislog

exit $retcode