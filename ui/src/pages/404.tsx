import { Button, Text, Center, Stack, Title } from "@mantine/core"
import { Link } from "@tanstack/react-router"

export const Error404 = () => {
  return (
    <Center h="100%">
      <Stack align="center" gap={0}>
        <Text fw={600}>404</Text>
        <Title fw={600} className="mt-2">Page not found</Title>
        <Text c="gray" className="mt-6">Sorry we couldn't find the page you're looking for</Text>
        <Button className="mt-6">
          <Link to="/">
            Go back home
          </Link>
        </Button>
      </Stack>
    </Center>
  )
}
