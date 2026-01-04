import { Page, PageTitle } from "@/components/atoms/Page"
import { GeneratorForm } from "@/components/generator/GeneratorForm"

export const GeneratorCreate = () => {
  return (
    <Page>
      <PageTitle
        title="Create a new playlist generator"
      />

      <GeneratorForm />
    </Page>
  )
}
