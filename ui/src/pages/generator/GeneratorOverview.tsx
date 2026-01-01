import { LinkButton } from "@/components/atoms/LinkButton"
import { Page, PageTitle, Section, SectionTitle } from "@/components/atoms/Page"
import { Group } from "@mantine/core"
import { LuSparkles } from "react-icons/lu"

export const GeneratorOverview = () => {
  return (
    <Page>
      <Group justify="space-between">
        <PageTitle
          title="Generate new playlists"
          description="Create playlists from presets and fine-tune them before saving."
        />
        <LinkButton to={"/generator/create"} leftSection={<LuSparkles />} radius="lg">New Generator</LinkButton>
      </Group>

      <Section>
        <SectionTitle
          title="Generated playlists"
        />
      </Section>
    </Page>
  )
}
