import http from 'k6/http';
import { check } from 'k6';

export const options = {
  scenarios: {
    race_condition: {
      executor: 'per-vu-iterations',
      vus: 100,
      iterations: 1,
    },
  },
};

export default function () {
  const url = 'http://localhost:8080/book';

  const payload = JSON.stringify({
    event_id: 1,
    user_id: `user_${__VU}`
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const res = http.post(url, payload, params);
  check(res, {
    'is status 200': (r) => r.status === 200,
  });
}