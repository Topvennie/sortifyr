import { Page, PageTitle, Section } from "@/components/atoms/Page"
import { GeneratorForm } from "@/components/generator/GeneratorForm"

export const GeneratorCreate = () => {
  return (
    <Page>
      <PageTitle
        title="Create a new playlist generator"
      />

      <Section>
        <GeneratorForm />
      </Section>
    </Page>
  )
}
