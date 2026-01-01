import { Page, PageTitle } from "@/components/atoms/Page"
import { useParams } from "@tanstack/react-router"

export const GeneratorEdit = () => {
  const { generatorId } = useParams({ from: "/public-layout/generator/edit/$generatorId" })

  return (
    <Page>
      <PageTitle
        title="Edit a playlist generator"
      />
    </Page>
  )
}
