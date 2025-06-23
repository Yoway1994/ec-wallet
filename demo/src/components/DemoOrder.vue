<script setup>
import { ref, reactive, computed, onUnmounted } from 'vue';

// State for form inputs
const price = ref('');
const orderId = ref('');
const selectedChain = ref('BNB'); // Fixed to BNB only
const selectedToken = ref('USDT'); // Default to USDT
const availableTokens = ['USDT', 'USDC'];

// API configuration
const baseUrl = 'http://localhost:8080';

// API response data
const paymentData = reactive({
  loading: false,
  error: null,
  data: null
});

// Generate a random order ID if none provided
const generateOrderId = () => {
  return 'ORD' + Date.now().toString().slice(-8) + Math.floor(Math.random() * 1000).toString().padStart(3, '0');
};

// Time remaining calculation
const timeRemaining = computed(() => {
  if (!paymentData.data || !paymentData.data.expire_time) return '';

  const now = new Date();
  const expireTime = new Date(paymentData.data.expire_time);
  const diffMs = expireTime - now;

  if (diffMs <= 0) return 'Expired';

  const minutes = Math.floor(diffMs / 60000);
  const seconds = Math.floor((diffMs % 60000) / 1000);

  return `${minutes}:${seconds.toString().padStart(2, '0')}`;
});

// Call API to generate payment address
const generatePaymentAddress = async () => {
  // Reset previous data
  paymentData.data = null;
  paymentData.error = null;
  paymentData.loading = true;

  try {
    // Use provided order ID or generate one
    const finalOrderId = orderId.value.trim() || generateOrderId();

    const response = await fetch(`${baseUrl}/api/v1/payment-orders`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        order_id: finalOrderId,
        chain: selectedChain.value,
        amount_usd: parseFloat(price.value) || 0,
        token: selectedToken.value
      })
    });

    if (!response.ok) {
      throw new Error(`API error: ${response.status}`);
    }

    const data = await response.json();
    paymentData.data = data;

    // Start countdown timer
    startCountdown();
  } catch (error) {
    console.error('Failed to generate payment address:', error);
    paymentData.error = error.message || 'Failed to generate payment address';
  } finally {
    paymentData.loading = false;
  }
};

// Timer for countdown
let countdownTimer = null;

const startCountdown = () => {
  // Clear any existing timer
  if (countdownTimer) {
    clearInterval(countdownTimer);
  }

  // Update every second
  countdownTimer = setInterval(() => {
    if (timeRemaining.value === 'Expired') {
      clearInterval(countdownTimer);
    }
  }, 1000);
};

// Clean up timer when component is unmounted
onUnmounted(() => {
  if (countdownTimer) {
    clearInterval(countdownTimer);
  }
});
</script>

<template>
  <div class="demo-order">
    <h1>Payment QR Code Generator</h1>

    <div class="form-container">
      <div class="form-group">
        <label for="price">Price Amount (USD):</label>
        <input
          id="price"
          v-model="price"
          type="number"
          placeholder="Enter USD amount"
          min="0"
        >
      </div>

      <div class="form-group">
        <label for="orderId">Order ID (optional):</label>
        <input
          id="orderId"
          v-model="orderId"
          type="text"
          placeholder="Leave empty for auto-generated ID"
        >
      </div>

      <div class="form-group">
        <label for="token">Payment Token:</label>
        <div class="token-selector">
          <div
            v-for="token in availableTokens"
            :key="token"
            class="token-option"
            :class="{ 'selected': selectedToken === token }"
            @click="selectedToken = token"
          >
            {{ token }}
          </div>
        </div>
      </div>

      <button
        @click="generatePaymentAddress"
        :disabled="paymentData.loading"
        class="generate-btn"
      >
        {{ paymentData.loading ? 'Generating...' : 'Generate Payment QR' }}
      </button>
    </div>

    <!-- Display payment data -->
    <div v-if="paymentData.data" class="payment-result">
      <h2>Payment Details</h2>

      <div class="result-card">
        <div class="order-info">
          <p><strong>Order ID:</strong> {{ paymentData.data.order_id }}</p>
          <p><strong>Chain:</strong> {{ paymentData.data.chain || selectedChain }}</p>
          <p><strong>Token:</strong> {{ paymentData.data.token || selectedToken }}</p>
          <p><strong>Address:</strong> {{ paymentData.data.address }}</p>
          <p><strong>Expires in:</strong> {{ timeRemaining }}</p>
        </div>

        <!-- QR Code placeholder - in a real app, use a QR code library here -->
        <div class="qr-code-container">
          <p>QR Code for:</p>
          <pre>{{ paymentData.data.address }}</pre>
          <!-- In a real app, replace with: <qrcode :value="paymentData.data.address" :options="{ width: 200 }" /> -->
        </div>
      </div>
    </div>

    <!-- Error message -->
    <div v-if="paymentData.error" class="error-message">
      <p>{{ paymentData.error }}</p>
    </div>
  </div>
</template>

<style scoped>
.token-selector {
  display: flex;
  gap: 1rem;
  margin-top: 0.5rem;
}

.token-option {
  flex: 1;
  padding: 0.75rem;
  text-align: center;
  border: 1px solid #ced4da;
  border-radius: 4px;
  cursor: pointer;
  background-color: #f8f9fa;
  color: #333333;
  transition: all 0.2s;
}

.token-option.selected {
  background-color: #4CAF50;
  color: white;
  border-color: #4CAF50;
  font-weight: bold;
}
.demo-order {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
  font-family: Arial, sans-serif;
  color: #333333;
}

h1 {
  text-align: center;
  margin-bottom: 2rem;
  color: #4CAF50
}

.form-container {
  background-color: #f8f9fa;
  padding: 1.5rem;
  border-radius: 8px;
  margin-bottom: 2rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.form-group {
  margin-bottom: 1rem;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: bold;
  color: #222222;
}

input, select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 1rem;
  color: #333333;
}

.generate-btn {
  background-color: #4CAF50;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 1rem;
  width: 100%;
  transition: background-color 0.3s;
  font-weight: bold;
}

.generate-btn:hover {
  background-color: #45a049;
}

.generate-btn:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.payment-result {
  background-color: #f8f9fa;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.payment-result h2 {
  color: #1a1a1a;
  margin-bottom: 1rem;
}

.result-card {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

@media (min-width: 768px) {
  .result-card {
    flex-direction: row;
  }

  .order-info, .qr-code-container {
    flex: 1;
  }
}

.order-info p {
  margin: 0.5rem 0;
  color: #222222;
}

.order-info strong {
  color: #000000;
}

.qr-code-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 1rem;
  background-color: white;
  border-radius: 4px;
  color: #222222;
}

pre {
  background-color: #f1f1f1;
  padding: 0.75rem;
  border-radius: 4px;
  overflow-x: auto;
  max-width: 100%;
  word-break: break-all;
  color: #000000;
}

.error-message {
  background-color: #f8d7da;
  color: #721c24;
  padding: 1rem;
  border-radius: 4px;
  margin-top: 1rem;
  font-weight: bold;
}
</style>
