<template v-slot="{dark}">
  <div id="app" class="container">
    <vue-particles
      color="#a9a9a9"
      :particleOpacity="0.7"
      :particlesNumber="80"
      shapeType="triangle"
      :particleSize="4"
      linesColor="#a9a9a9"
      :linesWidth="1"
      :lineLinked="true"
      :lineOpacity="0.6"
      :linesDistance="150"
      :moveSpeed="3"
      :hoverEffect="true"
      hoverMode="grab"
      :clickEffect="true"
      clickMode="push"
    >
    </vue-particles>
    <div class="row">
      <div class="col-md-6 offset-md-3 py-5">
        <h1>SHRTN⚡️</h1>
        <h3>Get a short address for any URL.</h3>
        <form v-on:submit.prevent="makeWebsiteThumbnail">
          <div class="form-group">
            <input
              v-model="websiteUrl"
              type="text"
              id="website-input"
              placeholder="Enter a website"
              class="form-control"
            />
          </div>
          <div class="form-group">
            <button class="btn btn-primary">Shorten!</button>
          </div>
        </form>
        <div id="result" class="container">
          <span style="font-size: 150%">
            <a :href="shortenedURL">{{ shortenedURL }} </a>
          </span>
          <button
            class="btn btn-primary"
            v-show="showResult"
            v-clipboard:copy="shortenedURL"
          >
            Copy
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";
export default {
  name: "App",

  data() {
    return {
      websiteUrl: "",
      shortenedURL: "",
      showResult: false,
    };
  },

  methods: {
    makeWebsiteThumbnail() {
      if (this.websiteUrl === "") {
        return;
      }
      this.shortenedURL = "";
      this.showResult = false;

      console.log(`shortening ${this.websiteUrl}`);
      axios
        .post("http://sho.rt:7000/url", {
          url: this.websiteUrl,
        })
        .then((response) => {
          this.shortenedURL = "http://sho.rt:7000/r/" + response.data.slug;
          this.showResult = true;
        })
        .catch((error) => {
          window.alert(`The API returned an error: ${error}`);
        });
    },
  },
};
</script>

<style>
:root {
  --bg: #fff;
  --color: #333333;
}

html.dark-mode {
  --bg: #232b32;
  --color: #ddd8ca;
}

body {
  background-color: var(--bg);
  color: var(--color);
}
</style>
