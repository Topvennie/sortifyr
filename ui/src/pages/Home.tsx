import { LoadingSpinner } from "@/components/molecules/LoadingSpinner"
import { useSettingGet } from "@/lib/api/setting"
import { useSync } from "@/lib/api/user"
import { formatDate } from "@/lib/utils"
import { Button, Title } from "@mantine/core"
import { notifications } from "@mantine/notifications"
import { FaCheck } from "react-icons/fa6"

export const Home = () => {
  const { data: settings, isLoading } = useSettingGet()

  const sync = useSync()

  const handleSync = () => {
    const id = notifications.show({
      loading: true,
      title: "Sync",
      message: "Synchronizing data with spotify",
      autoClose: false,
      withCloseButton: false,
    })

    sync.mutate(undefined, {
      onSuccess: () => notifications.update({
        id,
        variant: "succes",
        title: "Sync",
        message: "Synchronize done",
        icon: <FaCheck />,
        loading: false,
        autoClose: 3000,
      }),
    })
  }

  if (isLoading) return <LoadingSpinner />

  console.log(settings)

  return (
    <div className="flex flex-col justify-center items-center h-full pt-[10%]">
      <Title order={1} className="text-balance font-semibold tracking-tight">
        Welcome
      </Title>
      <p className="mt-6 text-pretty text-lg font-medium text-gray-500">
        {`Last sync: ${settings?.lastUpdate ? formatDate(settings?.lastUpdate) : "Never"}`}
      </p>
      <div className="mt-10">
        <Button onClick={handleSync}>
          Synchronize
        </Button>
      </div>
    </div>
  )
}
