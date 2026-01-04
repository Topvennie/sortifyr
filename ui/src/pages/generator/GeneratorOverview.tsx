import { LinkButton } from "@/components/atoms/LinkButton"
import { Page, PageTitle, Section, SectionTitle } from "@/components/atoms/Page"
import { Table } from "@/components/molecules/Table"
import { useGeneratorGetAll } from "@/lib/api/generator"
import { Generator } from "@/lib/types/generator"
import { ActionIcon, Badge, Group } from "@mantine/core"
import { useNavigate } from "@tanstack/react-router"
import { LuCheck, LuPencil, LuSparkles, LuTrash2, LuUndo2 } from "react-icons/lu"

export const GeneratorOverview = () => {
  const { data: generators, isLoading } = useGeneratorGetAll()

  const navigate = useNavigate()

  const handleEdit = (gen: Generator) => {
    navigate({ to: "/generator/edit/$generatorId", params: { generatorId: gen.id.toString() } })
  }

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

        <Table
          columns={[
            {
              accessor: "name",
              title: "Name & Description",
              render: ({ name, description }) => (
                <div>
                  <p className="font-semibold">{name}</p>
                  <p className="text-muted text-sm">{description}</p>
                </div>
              )
            },
            {
              accessor: "playlist",
              title: "Playlist",
              textAlign: "right",
              render: ({ playlistId }) => playlistId && <LuCheck className="ml-auto text-green-500 size-6" />,
            },
            {
              accessor: "maintained",
              title: "Maintained",
              textAlign: "right",
              render: ({ playlistId, maintained, intervalS }) => {
                if (!playlistId) return null
                const days = Math.floor(intervalS / (60 * 60 * 24))

                return <Badge color={maintained ? "secondary.1" : "gray"} className="ml-auto">{maintained ? `Every${days !== 1 ? ' ' + days : ''} day${days !== 1 ? 's' : ''}` : "One off"}</Badge>
              },
            },
            {
              accessor: "outdated",
              title: "Status",
              textAlign: "right",
              render: ({ maintained, outdated }) => maintained && <Badge color={outdated ? "red" : "gray"} className="ml-auto">{outdated ? "Outdated" : "Up to date"}</Badge>,
            },
            {
              accessor: "actions",
              title: "",
              width: 106,
              render: gen => (
                <div className="flex gap-0 flex-nowrap">
                  <div className="flex-1" />
                  <ActionIcon variant="subtle" color="black"><LuUndo2 /></ActionIcon>
                  <ActionIcon onClick={() => handleEdit(gen)} variant="subtle" color="black"><LuPencil /></ActionIcon>
                  <ActionIcon variant="subtle" color="red"><LuTrash2 /></ActionIcon>
                </div>
              )
            },
          ]}
          records={generators ?? []}
          fetching={isLoading}
          noRecordsText="No generators"
        />
      </Section>
    </Page>
  )
}
