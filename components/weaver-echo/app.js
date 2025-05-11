const express = require('express');
const bodyParser = require('body-parser');

const app = express();
const port = process.env.PORT || 8080;

app.use(bodyParser.json());

app.get('/health', (req, res) => {
  res.status(200).send('OK');
});

app.post('/', (req, res) => {
  const event = req.body;
  console.log('Echo received:', event.data);
  
  // No further re-emit
  res.status(204).send();
});

app.listen(port, () => {
  console.log(`Echo service listening on port ${port}`);
});
