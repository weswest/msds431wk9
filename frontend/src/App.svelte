<script>
  import {Greet} from '../wailsjs/go/main/App.js'
  import {CheckTermEmbed} from '../wailsjs/go/main/App.js'
  import {CheckDir} from '../wailsjs/go/main/App.js'
  import {RunDebug} from '../wailsjs/go/main/App.js'
  // @ts-ignore
  import heroImage from './assets/images/hero-image.png'

  let inputText = "";
  let outputText = "";

  const sample1 = `golang`; 
  const sample2 = `blueberry`; 

  function evaluateText() {
    console.log("evaluateText function called with input:", inputText);
    CheckTermEmbed(inputText).then(result => {
        console.log("Received result from backend:", result);
        outputText = result;
    });
  }

  function runDebug() {
    console.log("evaluateText function called with input:", inputText);
    RunDebug().then(result => {
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
              MVP.  Wk9 Assignment: Access Database
          </p>
          <h2>How to Use</h2>
          <p>
              Enter text and click evaluate
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
  <div class="button-container">
      <button on:click={evaluateText}>Evaluate</button>
      <button on:click={runDebug}>Debug</button>
  </div>
  <div class="output">{outputText}</div>
  <div class="rules">
      <h3>Sample bottom text</h3>
  </div>
</div>
