import { PlaylistTableView } from "@/components/playlist/PlaylistTableView"
import { PlaylistTreeView } from "@/components/playlist/PlaylistTreeView"
import { Switch, Title } from "@mantine/core"
import { useState } from "react"

export const Playlists = () => {
  const [treeView, setTreeView] = useState(false)

  return (
    <div className="flex flex-col gap-8">
      <div className="flex items-center justify-between">
        <Title order={1}>Playlists</Title>
        <Switch
          checked={treeView}
          onChange={e => setTreeView(e.target.checked)}
          label="Tree view"
        />
      </div>
      {treeView
        ? <PlaylistTreeView />
        : <PlaylistTableView />
      }
    </div>
  )
}

