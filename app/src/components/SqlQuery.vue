<script setup lang="ts">
import { useStore } from '@/store';
import { ref } from 'vue';

const query = ref("")
const queryResults = ref("")

const {sqlQuery} = useStore()

async function sendQuery() {
  const response = await sqlQuery(query.value)

  queryResults.value = JSON.stringify(response)
}
</script>
<template>
  <div class="container-fluid">
    <details>
      <summary>Advanced</summary>
      <article>
        <header>
          SQL Query Interface
        </header>
        <label for="query">Query</label>
        <input type="text" id="query" v-model="query" placeholder="SELECT * FROM ...">
        <button type="button" @click="sendQuery()">Submit Query</button>
        <label for="results">Results</label>
        <textarea id="results" v-model="queryResults"></textarea>
      </article>
    </details>
  </div>
</template>