<script setup lang="ts">
import { useStore } from '@/store';
import { onBeforeUnmount, ref } from 'vue';

const newOrderWidget = ref(0)
const newOrderQuantity = ref(0)

const { orders, widgets, deleteOrder, createOrder } = useStore()

function remove(ID: number) {
  deleteOrder(ID)
}

function create() {
  createOrder(+newOrderWidget.value, +newOrderQuantity.value)
}
</script>

<template>
  <div class="container-fluid">
    <article>
      <header>Orders</header>
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Widget</th>
            <th>Quantity</th>
            <th>&nbsp;</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="w in orders" v-bind:key="w.ID">
            <td>{{ w.ID }}</td>
            <td>{{ widgets.find(ww => ww.ID == w.Widget)?.Name }}</td>
            <td>{{ w.Quantity }}</td>
            <td><a class="danger" href="javascript://" @click="remove(w.ID)">Delete</a></td>
          </tr>
        </tbody>
      </table>
      <footer>
        <form>
          <div class="grid">
            <div>
              <select id="widget" required v-model="newOrderWidget">
                <option value="0">Select a widget...</option>
                <option v-for="w in widgets" v-bind:key="w.ID" :value="w.ID">{{  w.Name  }}</option>
              </select>
            </div>
            <div>
              <input type="text" v-model="newOrderQuantity" id="quantity" placeholder="Quantity" />
            </div>
            <div>
              <button type="button" @click="create()" :disabled="newOrderWidget == 0 || newOrderQuantity == 0">Create Order</button>
            </div>
          </div>
        </form>
      </footer>
    </article>
  </div>
</template>