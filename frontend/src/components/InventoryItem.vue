<template>
    <div class="grid container mx-1 p-0">
        <div class="col-9 flex justify-content-center align-items-center">
            {{ name }}
        </div>
        <div class="col-3 flex justify-content-center align-items-center quantity-container">
            {{ quantity }} {{ unit }}
        </div>
    </div>
</template>

<script setup>
defineProps({
    name: { type: String, default: '' },
    quantity: { type: [String, Number], default: 0 },
    unit: { type: String, default: '' },
})
</script>

<style scoped>
.container {
    border: 1px solid #ccc;
    border-radius: 5px;
    padding: 5px;
    margin:0px;
}
.quantity-container{
    border-top-right-radius:5px;
    border-bottom-right-radius:5px;
    background-color:bisque;
    height:3rem;
}
</style>
