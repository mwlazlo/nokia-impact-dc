curl -v -X PUT \
 --header 'Content-Type: application/json' --header 'Accept: application/json' \
--header 'Authorization: Basic VWJpaWsuVEg6VWJpaWtAMTk=' -d  '
{ "headers": {"authorization":"Basic dWF0YWRlcDpBc2RmMSM="}, "url":
 "http://localhost:8080/m2m/impact/callback"
 }' https://impact.idc.nokia.com/m2m/applications/registration
