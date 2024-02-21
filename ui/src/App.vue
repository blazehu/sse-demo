<template>
  <div id="app">
    <img alt="Vue logo" src="./assets/logo.png">
    <HelloWorld :msg="msg"/>
  </div>
</template>

<script>
import HelloWorld from './components/HelloWorld.vue';

export default {
  name: 'App',
  components: {
    HelloWorld
  },
  data() {
    return {
      msg: '',
    };
  },
  mounted() {
    this.loadData();
  },
  methods: {
    async loadData() {
      const vm = this;
      const eventSource = new EventSource("http://127.0.0.1:8080/api/v1/stream/chat");
      eventSource.onopen = function () {
        console.log('connect eventSource success.');
      };
      eventSource.onmessage =  (e) => {
        const data = JSON.parse(e.data);
        vm.msg = `${data.result.content} - ${data.result.user} send`;
      };
      eventSource.onerror = function () {
        console.log('connect eventSource failed.');
      };
      this.$forceUpdate();
    },
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
