//import k6 and do a request to localhost/events with a POST method passing the event object as a JSON string
import http from 'k6/http';
import { sleep, check } from 'k6';


export function teardown() {
    //delte all events created by the load test
    const url = 'http://localhost:8000/event';
    const params = {
        headers: {
            'Content-Type': 'application/json',
            'accept': '*/*',
        },
    }
    let res = http.del(url, null, params)
    check(res, {
        'is status 200': (r) => r.status === 200,
    });
}

export default function () {
    const url = 'http://localhost:8000/event';
    const params = {
        headers: {
            'Content-Type': 'application/json',
            'accept': '*/*',
        },
    }
    const payload = JSON.stringify({
        "title": "string",
        "description": "string",
        "location": "string",
        "start_time": "2023-05-07T13:22:08.124Z",
        "end_time": "2023-05-07T13:22:08.124Z",
        "instagram_page": "string"
    })
    let res = http.post(url, payload, params)
    check(res, {
        'is status 200': (r) => r.status === 200,
    });
}