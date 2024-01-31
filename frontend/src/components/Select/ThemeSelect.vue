<template>
    <!-- e -->
    <span v-if="options != undefined && options != null && options.length != 0" class="ThemeSelect">
        <!-- a -->
        <span class="custom-select" @blur="open = false">
            <!-- b -->
            <div class="selected" :class="{ open: open }" @click="open = !open">
                <span v-if="modelValue == ''">Cyan</span>
                <span v-else> {{ selected.name }}</span>
            </div>
            <span class="items" :class="{ selectHide: !open }">
                <div
                    v-for="(option, i) of options"
                    :key="i"
                    @click="
                        () => {
                            selected = option
                            modelValue = option.id
                            open = false
                            console.log(option.id)
                            setTheme(option.id)
                        }
                    ">
                    <span>
                        <span>
                            {{ option.name }}
                        </span>
                    </span>
                </div>
            </span>
        </span>
    </span>
</template>

<script lang="ts">
import { defineComponent } from "vue"

interface Theme {
    id: string
    name: string
}

const themes: Array<Theme> = [
    {
        id: "",
        name: "Cyan"
    },
    {
        id: "blue",
        name: "Blue"
    },
    {
        id: "green",
        name: "Green"
    },
    {
        id: "orange",
        name: "Orange"
    },
    {
        id: "purple",
        name: "Purple"
    },
    {
        id: "red",
        name: "Red"
    },
    {
        id: "yellow",
        name: "Yellow"
    },
    {
        id: "black",
        name: "dark"
    }
]

export default defineComponent({
    name: "ThemeSelect",
    data() {
        return {
            modelValue: "",
            selected: {} as Theme,
            open: false,
            options: themes
        }
    },
    watch: {
        modelValue: function () {
            this.options.forEach(element => {
                if (element.id == this.modelValue) {
                    this.selected = element
                }
            })
        }
        // options: function () {
        //     this.options.forEach(element => {
        //         if (element.id == this.modelValue) {
        //             this.selected = element
        //         }
        //     })
        // }
    },
    mounted() {
        if (localStorage.theme) {
            this.modelValue = localStorage.theme
            this.setTheme(localStorage.theme)
        }
        // this.options.forEach(element => {
        //     if (element.id == this.modelValue) {
        //         this.selected = element
        //     }
        // })
    },
    methods: {
        useSelected: function () {
            this.options.forEach(element => {
                if (element.id == this.modelValue) {
                    this.selected = element
                    console.log(element.id)
                    document.firstElementChild?.setAttribute("data-theme", element.id)
                    localStorage.theme = element.id
                }
            })
        },
        setTheme: function (theme: string) {
            document.firstElementChild?.setAttribute("data-theme", theme)
            localStorage.theme = theme
        }
    }
})
</script>

<style lang="scss" scoped>
img {
    height: 12px;
    padding-right: 5px;
}

.custom-select {
    // position: relative;
    text-align: left;
    outline: none;
    width: max-content;
}

.custom-select .selected {
    background-color: var(--color);
    border-radius: 6px;
    border: 1px solid transparent;
    color: var(--font-color);
    padding: 5px;
    cursor: pointer;
    padding-right: 25px;
    width: max-content;

    user-select: none;
}

.custom-select .selected.open {
    border: 1px solid var(--color-half);
    border-radius: 6px 6px 0px 0px;
}

.custom-select .selected:after {
    position: absolute;
    content: "";
    top: 50%;
    // top: 22px;
    // right: 5px;
    // left: 5px;
    width: 0;
    height: 0;
    border: 5px solid transparent;
    border-color: #fff transparent transparent transparent;
}

.custom-select .items {
    color: var(--font-color);
    border-radius: 0px 0px 6px 6px;
    overflow: hidden;
    border-right: 1px solid var(--color);
    border-left: 1px solid var(--color);
    border-bottom: 1px solid var(--color);
    position: absolute;
    min-width: min-content;
    width: max-content;
    background-color: var(--color-dark);
    left: 0;
    right: 0;
    z-index: 1;
}

.custom-select .items div {
    color: var(--font-color);
    // padding-left: 1em;
    padding: 5px;
    cursor: pointer;
    user-select: none;
}

.custom-select .items div:hover {
    background-color: var(--color-full);
}

.selectHide {
    display: none;
}

.custom-select {
    position: relative;
    display: block;
    width: min-content;
    height: min-content;
}
</style>
