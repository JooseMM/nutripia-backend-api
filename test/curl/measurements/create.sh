curl -v -X POST http://localhost:3000/body-measurements/$1 \
  -H "Content-Type: application/json" \
  -d '{
    "mass": 75.5,
    "stature": 160.2,
    "sittingHeight": 90.5,
    "armSpan": 162.0,
    "triceps": 12.4,
    "subscapular": 15.2,
    "biceps": 8.1,
    "iliacCrest": 14.5,
    "supraspinale": 10.2,
    "abdominal": 18.3,
    "frontThigh": 16.1,
    "medialCalf": 9.4,
    "head": 56.0,
    "neck": 38.5,
    "armRelaxed": 32.0,
    "armFlex": 35.5,
    "forearm": 28.2,
    "wrist": 17.5,
    "chest": 102.0,
    "waist": 85.5,
    "hip": 98.0,
    "thighHigh": 60.2,
    "thighLow": 54.1,
    "calf": 37.5,
    "ankle": 22.0
  }' | jq
