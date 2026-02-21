import { check } from 'k6';
import http from 'k6/http';

export const options = {
  scenarios: {
    thundering_herd: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '2s', target: 500 },
        { duration: '8s', target: 500 },
      ],
    },
  },
};

export default function () {
  const res = http.get('http://localhost:8080/event/1');
  check(res, {
    'is_status_200': (r) => r.status === 200,
  });
}