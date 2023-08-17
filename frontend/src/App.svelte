<script>
//  import {Greet} from '../wailsjs/go/main/App.js'
  import {CheckTermEmbed} from '../wailsjs/go/main/App.js'
  // @ts-ignore
  import heroImage from './assets/images/hero-image.png'

  let inputText = "";
  let outputText = "";

  const validTerms = [
    "break", "case", "chan", "const", "continue", "defer", "else", "fallthrough",
    "for", "func", "go", "goto", "if", "import", "interface", "map", "package",
    "range", "return", "select", "struct", "switch", "type", "var"
  ];

  const sillyTerms = [
    "blueberry", "dreams", "unicorn", "rainbow", "chocolate", "fairy", "moonbeam",
    "sparkle", "glitter", "bubble", "whimsy", "twinkle", "fizz", "flutter", "giggle"
  ];
  function evaluateText() {
    console.log("evaluateText function called with input:", inputText);
    CheckTermEmbed(inputText).then(result => {
        console.log("Received result from backend:", result);
        outputText = result;
    });
  }

  function loadSample(isValid) {
    if (isValid) {
      inputText = validTerms[Math.floor(Math.random() * validTerms.length)];
    } else {
      inputText = sillyTerms[Math.floor(Math.random() * sillyTerms.length)];
    }
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
          <h2>Week 9 Assignment: Access A Database</h2>
          <p>
              This assignment mimicks a desktop application that allows the user to ask questions of an LLM and receive a response.  The functionality of posing a question to an LLM and receiving an answer has been replaced with a database lookup.
          </p>
          <h2>How to Use</h2>
          <p>
              Enter text and click evaluate.  If the term exists in the database, the definition will be returned.  If the term does not exist in the database, an error message will be returned.  Note that the database search is case sensitive and requires an exact match.
          </p>
          <p>
            The "Valid Term" button on the right will randomly select a term that exists in the database.  Note this has been hard-coded.
          </p>
      </div>
  </div>
  <div class="input-section">
      <textarea bind:value={inputText} placeholder="Enter your text here..."></textarea>
      <div class="sample-buttons">
          <button on:click={() => loadSample(true)}>Valid Term</button>
          <button on:click={() => loadSample(false)}>Invalid Term</button>
      </div>
  </div>
  <div class="button-container">
      <button on:click={evaluateText}>Evaluate</button>
  </div>
  <div class="output">{outputText}</div>
  <div class="rules">
      <h3>Valid Terms</h3>
      <p>
        "break", "case", "chan", "const", "continue", "defer", "else", "fallthrough",
        "for", "func", "go", "goto", "if", "import", "interface", "map", "package",
        "range", "return", "select", "struct", "switch", "type", "var"
    </p>
</div>
</div>
