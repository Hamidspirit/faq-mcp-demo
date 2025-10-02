require('dotenv').config();
const express = require('express');
const { Anthropic } = require('@anthropic-ai/sdk');
const fs = require('fs');

const app = express();
app.use(express.json());

const client = new Anthropic({ apiKey: process.env.ANTHROPIC_API_KEY });

// Load FAQ data
const faqs = JSON.parse(fs.readFileSync('./data/faqs.json', 'utf8'));

// Simple keyword match for now (we’ll add embeddings later)
function findBestMatch(question) {
  let best = faqs[0];
  let score = 0;
  for (let faq of faqs) {
    let overlap = 0;
    const qWords = question.toLowerCase().split(" ");
    for (let w of qWords) {
      if (faq.question.toLowerCase().includes(w)) overlap++;
    }
    if (overlap > score) {
      score = overlap;
      best = faq;
    }
  }
  return best;
}

app.post('/api/chat', async (req, res) => {
  const { question } = req.body;
  const match = findBestMatch(question);

  // Call Anthropic to rephrase + polish answer
  const response = await client.messages.create({
    model: 'claude-3-5-sonnet-20240620',
    max_tokens: 300,
    messages: [
      {
        role: 'system',
        content: "You are a professional FAQ assistant. Use the provided answer, but make it clear and helpful."
      },
      {
        role: 'user',
        content: `Q: ${question}\nSuggested answer: ${match.answer}`
      }
    ]
  });

  res.json({ answer: response.content[0].text });
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(`FAQ bot listening at http://localhost:${PORT}`);
});
