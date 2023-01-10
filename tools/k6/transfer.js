import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '5s', target: 101 }, // slow traffic
    ],
    thresholds: {
        http_req_duration: ['med<7000'],
    },
};


export default function () {
    const url = 'http://localhost:3000/transfer'
    const response = http.post(url,
        JSON.stringify({
            "from": 1,
            "to": 2,
            "amount": 1
        }),
        {
            headers: {
                "Content-Type": "application/json"
            }
        }
    );

    if (response.status !== 200) {
        console.log(`Request failed with status: ${response.status}`, response.json())
    }
    console.log("successfully transferred", response.json())
    check(response, { 'status was 200': ({status}) => status === 200 });
    sleep(1);
}
