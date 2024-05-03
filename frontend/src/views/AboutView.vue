<template>
    <div class="AboutView">
        <h2>About</h2>
        <ul>
            <li><strong>Name:</strong> Fabio Kälin</li>
            <li><strong>E-Mail:</strong> fabio.kaelin.fk+oauth@gmail.com</li>
            <li><strong>GitHub:</strong> <a href="https://github.com/FabioKaelin">FabioKaelin</a></li>
        </ul>

        <h3>Haftungsausschluss</h3>
        <p>
            Der Inhalt dieser Website wurde sorgfältig geprüft und zusammengestellt. Trotzdem kann keine Gewähr für die Richtigkeit, Vollständigkeit und Aktualität der bereitgestellten Informationen
            übernommen werden.
        </p>

        <h3>Urheberrecht</h3>
        <p>Alle Inhalte auf dieser Website werden ausschließlich zu Informationszwecken bereitgestellt. Du darfst die Inhalte dieser Website für private, nichtkommerzielle Zwecke nutzen.</p>

        <h3>Datenschutz</h3>
        <p>
            Die Nutzung dieser Website ist in der Regel ohne Angabe personenbezogener Daten möglich. Personenbezogene Daten (beispielsweise Name, Anschrift oder E-Mail-Adresse) werden auf dieser
            Website nur erhoben, soweit dies technisch erforderlich ist.
        </p>

        <h3>Kontakt</h3>
        <p>Bei Fragen oder Anregungen zu dieser Website oder den hier bereitgestellten Informationen kannst du mich gerne per E-Mail kontaktieren.</p>
    </div>
</template>

<script lang="ts">
import { defineComponent } from "vue"

export default defineComponent({
    name: "AboutView"
})
</script>

<style scoped lang="scss">
.AboutView {
    background-color: var(--color);
    border-radius: 10px;
    padding: 5px;
}

a{
    color: var(--color-full);
    // text-decoration: none;
}
</style>
