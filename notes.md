- https://github.com/ai/audio-recorder-polyfill
- https://github.com/ReactTraining/react-router/issues/7297
  (routing anims are not working atm)
- https://github.com/faroit/awesome-python-scientific-audio#read-write
- https://codesandbox.io/s/81zkxw8qnl?file=/src/index.tsx:151-240
- https://github.com/processing/p5.js-sound 
  (local processing)
- https://p1wmc.csb.app/ 
  (cool background css tricks)
- https://github.com/GoogleCloudPlatform/gke-managed-certs/issues/13
  (my dns issue with LB&ingress)
  also relevant: https://stackoverflow.com/questions/53886750/google-managed-ssl-certificate-stuck-on-failed-not-visible
  

- orms
  - go sql tools: https://github.com/Masterminds/squirrel with https://github.com/Masterminds/structable
  - ts orm: https://github.com/typeorm/typeorm#Installation
- auth
  - keycloak might be good to use after all? allows signin through other vendors, which could help adoption? https://www.keycloak.org/getting-started/getting-started-kube
  - authboss (golang)
  - https://docs.expo.io/versions/latest/sdk/auth-session/
  - jwt is good for auth from other vendors, e.g. keycloak, google etc., if the backend needs to handle the auth itself jwt is unnecessary? http://cryto.net/~joepie91/blog/2016/06/13/stop-using-jwt-for-sessions/
- infra
  - gh deploying cd with k8: https://docs.github.com/en/actions/guides/deploying-to-google-kubernetes-engine
- geo types
  - distances calc https://www.nhc.noaa.gov/gccalc.shtml