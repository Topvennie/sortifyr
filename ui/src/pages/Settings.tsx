import { Button } from "@/components/atoms/Button"
import { Page, PageTitle, Section, SectionTitle } from "@/components/atoms/Page"
import { useSettingUploadExport } from "@/lib/api/setting"
import { CONTENT_TYPE } from "@/lib/types/contentType"
import { getErrorMessage } from "@/lib/utils"
import { FileButton, Group, Pill, Stack } from "@mantine/core"
import { notifications } from "@mantine/notifications"
import { LuCloudUpload } from "react-icons/lu"

export const Settings = () => {
  const uploadExport = useSettingUploadExport()

  const handleExport = (file: File | null) => {
    if (!file) return

    uploadExport.mutateAsync(file, {
      onSuccess: () => notifications.show({ message: "Import started. Check the task page to see the progress." }),
      onError: async (error) => {
        const msg = await getErrorMessage(error)
        notifications.show({ color: "red", message: msg })
      }
    })
  }

  return (
    <Page>
      <PageTitle
        title="Settings"
      />

      <Section className="flex-none">
        <SectionTitle
          title="General"
        />

        <Stack gap="md">
          <Group>
            <p className="font-bold">Import your spotify data</p>
            <Pill bg="secondary.1">Important</Pill>
          </Group>
          <Stack gap={0}>
            <p className="mb-2">Sortfiyr works best with a full export of your playlists and listening history. You can request a copy of your data from Spotify and then import it here.</p>
            <p>1. Go to your <a href="https://www.spotify.com/account/overview/" target="_blank" rel="noopener noreferrer">account page</a>.</p>
            <p>2. Go to account privacy and request extended streaming history.</p>
            <p>3. When Spotify emails your the download link, save the .zip file to your computer.</p>
            <p>4. Click the button below and choose the .zip file to start the import.</p>
          </Stack>
          <p className="font-semibold text-red-500">This will clear all historic data before the most recent song in the zip file.</p>
          <FileButton onChange={handleExport} accept={CONTENT_TYPE.ZIP}>
            {(props) => (
              <div>
                <Button leftSection={<LuCloudUpload />} {...props}>Upload Spotify data</Button>
              </div>
            )}
          </FileButton>
        </Stack>
      </Section>
    </Page>
  )
}
