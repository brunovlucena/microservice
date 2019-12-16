import http from "k6/http";
import { group, check  } from "k6";

export let options = {
    vus: 30,
    thresholds: {
        //'http_req_duration{kind:html}': ["avg<=10"],
        'http_reqs': ["rate>100"],
        "http_req_duration": ["p(95)<500"],
        "check_failure_rate": [
            // Global failure rate should be less than 1%
            "rate<0.01",
            // Abort the test early if it climbs over 5%
            { threshold: "rate<=0.05", abortOnFail: true },
        ],
    },
};

export default function() {
    group("GET /configs", function() {
        let res = http.get("http://127.0.0.1:8000/configs", { headers: { "Content-Type": "application/json" }});
        check(res, {
            "status is 302": (r) => r.status === 302,
        });
    });

    group("GET /config/pod-2", function() {
        let res = http.get("http://127.0.0.1:8000/configs/pod-2", { headers: { "Content-Type": "application/json" }});
        check(res, {
            "status is 302": (r) => r.status === 302,
        });
    });

    let payload = '{"metadata":"monitoring":{"enabled":false}},"name":"pod-new"}';
    let body = JSON.stringify(payload);
    group("POST /configs", function() {
        let res = http.post("http://127.0.0.1:8000/configs", { verb: "post" }, body, { headers: { "Content-Type": "application/json" }});

        let j = JSON.parse(res.body);

        check(res, {
            "status is 201": (r) => r.status === 201,
        });
    });

    let payload2 = '{"metadata":"monitoring":{"enabled":true}},"name":"pod-new"}';
    let body2 = JSON.stringify(payload2);
    group("PUT /configs/pod-3", function() {
        let res = http.put("http://127.0.0.1:8000/configs/pod-3", { verb: "put" }, body2, { headers: { "Content-Type": "application/json" }});
        check(res, {
            "status is 200": (r) => r.status === 302,
        });
    });

    group("DELETE /configs/pod-1-idonotexist", function() {
        let res = http.del("http://127.0.0.1:8000/configs/pod-1-idonotexist", { headers: { "Content-Type": "application/json" }});
        check(res, {
            "status is 422": (r) => r.status === 422,
        });
        let res2 = http.del("http://127.0.0.1:8000/configs/pod-new");
        check(res2, {
            "status is 200": (r) => r.status === 200,
        });
    });
};
