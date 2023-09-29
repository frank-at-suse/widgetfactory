<script setup lang="ts">
import WidgetTable from './components/WidgetTable.vue';
import OrderTable from './components/OrderTable.vue';
import SqlQuery from './components/SqlQuery.vue';
</script>

<template>
  <main class="container-fluid">
    <h1>Widget Factory</h1>
    <div class="grid">
      <div>
        <WidgetTable />
      </div>
      <div>
        <OrderTable />
      </div>
    </div>
    <div class="grid">
      <SqlQuery />
    </div>
  </main>
</template>