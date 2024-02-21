<template>
    <div class="dsg">
        <h3>Datenschutzerklärung</h3>

        <p>
            Diese Datenschutzerklärung gilt für alle Benutzer innerhalb der Europäischen Union (EU) sowie der Schweiz und berücksichtigt die Anforderungen der Datenschutz-Grundverordnung (DSGVO) der
            EU und der Datenschutzgesetze der Schweiz.
        </p>

        <h4>1. Verantwortlichkeit für die Datenverarbeitung:</h4>
        <p>Wir sind verantwortlich für die Verarbeitung Ihrer personenbezogenen Daten gemäß den Bestimmungen der DSGVO und den geltenden Datenschutzgesetzen der Schweiz.</p>

        <h4>2. Zweck der Datenverarbeitung:</h4>
        <p>Wir verarbeiten personenbezogene Daten ausschließlich für klar definierte und rechtmäßige Zwecke und beschränken die Verarbeitung auf das für die Zwecke erforderliche Maß.</p>

        <h4>3. Datenerhebung und Verwendung:</h4>
        <ul>
            <li>Wir erheben personenbezogene Daten nur, wenn sie freiwillig von Ihnen zur Verfügung gestellt werden oder wenn wir gesetzlich dazu verpflichtet sind.</li>
            <li>
                Die von Ihnen bereitgestellten Daten können zur Verbesserung unserer Dienstleistungen, zur Personalisierung Ihrer Erfahrung oder zur Bearbeitung von Transaktionen verwendet werden.
            </li>
        </ul>

        <h4>4. Datensicherheit:</h4>
        <ul>
            <li>Wir ergreifen angemessene Sicherheitsmaßnahmen, um Ihre persönlichen Daten vor Verlust, Missbrauch, unbefugtem Zugriff oder Offenlegung zu schützen.</li>
            <li>Dennoch können wir die Sicherheit Ihrer Daten nicht vollständig garantieren, da keine Übertragung über das Internet oder elektronische Speicherung absolut sicher ist.</li>
        </ul>

        <h4>5. Rechte der Benutzer:</h4>
        <ul>
            <li>
                Gemäß der DSGVO und den Datenschutzgesetzen der Schweiz haben Benutzer das Recht auf Zugang, Berichtigung, Löschung und Widerspruch gegen die Verarbeitung ihrer personenbezogenen
                Daten.
            </li>
            <li>Wir werden angemessene Schritte unternehmen, um diese Rechte zu respektieren und umzusetzen.</li>
        </ul>

        <h4>6. Internationale Datenübertragungen:</h4>
        <ul>
            <li>
                Unter Umständen übertragen wir personenbezogene Daten außerhalb der EU und der Schweiz. In diesem Fall werden angemessene Schutzmaßnahmen getroffen, um sicherzustellen, dass Ihre Daten
                angemessen geschützt sind.
            </li>
        </ul>

        <h4>7. Kontakt:</h4>
        <p>Für Fragen oder Anliegen bezüglich Ihrer Datenschutzrechte oder dieser Datenschutzerklärung können Sie sich gerne an uns wenden.</p>

        <p>
            Diese Datenschutzerklärung kann von Zeit zu Zeit aktualisiert werden, um Änderungen in unseren Datenschutzpraktiken und in der geltenden Gesetzgebung widerzuspiegeln. Wir empfehlen Ihnen,
            diese Datenschutzerklärung regelmäßig zu überprüfen, um über unsere Datenschutzpraktiken informiert zu bleiben.
        </p>
        <br /><br />
    </div>
</template>

<script lang="ts">
import { defineComponent } from "vue"
// import ChampionGuess from '@/components/ChampionGuess.vue';

export default defineComponent({
    name: "DSGView"
})
</script>

<style scoped lang="scss">
.dsg {
    background-color: var(--color);
    border-radius: 10px;
    padding: 5px;
}
</style>
