<script>
  import {Greet} from '../wailsjs/go/main/App.js'
  // @ts-ignore
  import heroImage from './assets/images/hero-image.png'

  let inputText = "";
  let outputText = "";

  const sample1 = `golang`; 
  const sample2 = `blueberry`; 

  function evaluateText() {
    console.log("evaluateText function called with input:", inputText);
    Greet(inputText).then(result => {
        console.log("Received result from backend:", result);
        outputText = result;
    });
  }

  function loadSample(sample) {
    inputText = sample;
  }


</script>

<style>
  .container {
      padding: 20px;
      display: flex;
      flex-direction: column;
      align-items: center;
  }

  .header-section {
    display: flex;
    align-items: center; 
    gap: 20px; 
}


  .hero-image {
    max-width: 200px; 
    margin-right: 20px;
  }

  .input-section {
      display: flex;
      align-items: start;
      width: 90%;
  }

  textarea {
      flex: 1;
      height: 200px;
      margin-right: 20px;
  }

  .sample-buttons {
      display: flex;
      flex-direction: column;
  }

  button {
      margin-bottom: 20px;
  }

  .output {
      width: 90%;
      border: 1px solid #ccc;
      padding: 10px;
      white-space: pre-line;
      text-align: left; 
  }

  .rules {
      text-align: left;
      font-size: 0.8rem;
  }
</style>

<div class="container">
  <div class="header-section">
      <!-- svelte-ignore a11y-img-redundant-alt -->
      <img src={heroImage} alt="Hero Image" class="hero-image" />
      <div>
          <h2>Assignment:</h2>
          <p>
              Create a program using Wails, Svelte, and Vale that allows users to enter text and receive feedback on their grammar. As a test case, focus is on whether "data" should be plural or singular.
          </p>
          <h2>How to Use</h2>
          <p>
              Enter text and click evaluate. The only rules turned on are the data singular / plural, as defined by the radio boxes below. Please see the bottom for the specific rules that are captured. Note that these rules are hard-coded: with more time we can add more sophistication.
          </p>
          <p>
              If you want some sample text to get you started, please click the samples on the right.
          </p>
      </div>
  </div>
  <div class="input-section">
      <textarea bind:value={inputText} placeholder="Enter your text here..."></textarea>
      <div class="sample-buttons">
          <button on:click={() => loadSample(sample1)}>Sample 1</button>
          <button on:click={() => loadSample(sample2)}>Sample 2</button>
      </div>
  </div>
  <button on:click={evaluateText}>Evaluate</button>
  <div class="output">{outputText}</div>
  <div class="rules">
      <h3>Sample bottom text</h3>

  </div>
</div>
