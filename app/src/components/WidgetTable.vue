<script setup lang="ts">
import { useStore } from '@/store';
import { ref } from 'vue';

const newWidgetName = ref("");

const { widgets, deleteWidget, createWidget } = useStore()

function remove(ID: number) {
  deleteWidget(ID)
}

function create() {
  createWidget(newWidgetName.value)
  newWidgetName.value = "";
}

</script>

<template>
  <div class="container-fluid">
    <article>
      <header>Widgets</header>
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>&nbsp;</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="w in widgets" v-bind:key="w.ID">
            <td>{{ w.ID }}</td>
            <td>{{ w.Name }}</td>
            <td><a class="danger" @click="remove(w.ID)" href="javascript://">Delete</a></td>
          </tr>
        </tbody>
      </table>
      <footer>
        <form>
          <div class="grid">
            <div>
              <input type="text" v-model="newWidgetName" placeholder="Widget name">
            </div>
            <div>
              <button :disabled="newWidgetName == ''" @click="create()" type="button">Create widget</button>
            </div>
          </div>
        </form>
      </footer>
    </article>
  </div>
</template>