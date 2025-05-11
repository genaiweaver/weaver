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
  console.log('Planner received:', event.data);
  
  const response = {
    data: {
      plannedFrom: event.data
    }
  };
  
  res.json(response);
});

app.listen(port, () => {
  console.log(`Planner service listening on port ${port}`);
});
