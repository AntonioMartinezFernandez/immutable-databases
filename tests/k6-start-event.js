import http from 'k6/http';
import { sleep } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

const SERVICE_URL = 'http://localhost:3000/api/tracking/events';

export const options = {
  stages: [
    {
      duration: '5s',
      target: 1000,
    },
    {
      duration: '20s',
      target: 1000,
    },
    {
      duration: '5s',
      target: 0,
    },
  ],
};

export default function () {
  const url = SERVICE_URL;
  const tenantId = 'bab1fc99-df84-4239-9998-957039d515b4';
  const transactionId = uuidv4();
  const payload = {
    version: 2,
    operationId: transactionId,
    tenantId,
    sessionId: transactionId,
    source: 'sdk.mobile',
    family: 'ONBOARDING',
  };

  const params = {
    headers: {
      Authorization: 'Basic bWFub2xpOmhvbGk=', // user:manoli password:holi
      'Content-Type': 'application/json',
      'x-api-key': '000',
    },
  };

  makeStartRequest(url, params, payload);
  sleep(1);
}

function makeStartRequest(url, params, data) {
  const stepTypes = {
    START: { value: 'START', stepId: uuidv4() },
    SELPHID_WIDGET: { value: 'SELPHID_WIDGET', stepId: uuidv4() },
    OPERATION_RESULT: { value: 'OPERATION_RESULT', stepId: uuidv4() },
  };

  const payload = JSON.stringify(
    Object.assign(
      {
        events: [
          {
            eventId: uuidv4(),
            clientTimestamp: new Date().toISOString(),
            executionTime: null,
            payload: {
              type: 'STEP_CHANGE',
              stepType: stepTypes.START.value,
              stepId: stepTypes.START.stepId,
              component: {
                id: 'tracking_component',
                version: '2.0.2-SNAPSHOT',
              },
            },
          },
        ],
      },
      data
    )
  );

  http.post(url, payload, params);
}
