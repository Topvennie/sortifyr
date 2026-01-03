import { Page, PageTitle } from "@/components/atoms/Page"
import { LoadingSpinner } from "@/components/molecules/LoadingSpinner"
import { useGeneratorGetAll } from "@/lib/api/generator"
import { useParams } from "@tanstack/react-router"
import { Error404 } from "../404"
import { GeneratorForm } from "@/components/generator/GeneratorForm"

export const GeneratorEdit = () => {
  const { generatorId } = useParams({ from: "/public-layout/generator/edit/$generatorId" })
  const { data: generators, isLoading } = useGeneratorGetAll()

  if (isLoading) return <LoadingSpinner />

  const generator = generators?.find(g => g.id.toString() === generatorId)
  if (!generator) return <Error404 />

  return (
    <Page>
      <PageTitle
        title="Edit a playlist generator"
      />

      <GeneratorForm generator={generator} />
    </Page>
  )
}
