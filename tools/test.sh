#!/bin/bash
#

echo "Start testing ..."
exit_status="0"
echo -e "\nTesting JSON values ..."

json1=$(./awesome-ci parseJSON -file ./tests/testdata/values.json -value .value1)
json2=$(./awesome-ci parseJSON -file ./tests/testdata/values.json -value .value2)
# Not yet implemented
json3=$(./awesome-ci parseJSON -file ./tests/testdata/values.json -value .deepObject.value1)

if [ "$json1" == "hello" ] && [ "$json2" == "world" ]; then 
    echo "JSON testing successful"
else
    echo "JSON testing not successful. Values: $json1 $json2 $json3"
    exit_status="1"
fi

echo -e "\nTesting YAML values ..."

json1=$(./awesome-ci parseYAML -file ./tests/testdata/values.yaml -value .value1)
json2=$(./awesome-ci parseYAML -file ./tests/testdata/values.yaml -value .value2)
# Not yet implemented
json3=$(./awesome-ci parseYAML -file ./tests/testdata/values.yaml -value .deepObject.value1)

if [ "$json1" == "hello" ] && [ "$json2" == "world" ]; then 
    echo "YAML testing successful"
else
    echo "YAML testing not successful. Values: $json1 $json2 $json3"
    exit_status="1"
fi

exit $exit_status