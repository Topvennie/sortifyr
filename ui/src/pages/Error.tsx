import { isResponseNot200Error } from "@/lib/api/query";
import { Button, Container, Title } from "@mantine/core";
import { ErrorComponentProps, useNavigate } from "@tanstack/react-router";
import { Error404 } from "./404";
import { useAuth } from "@/lib/hooks/useAuth";

export const Error = ({ error, reset }: ErrorComponentProps) => {
  const navigate = useNavigate()
  const { logout } = useAuth()

  if (isResponseNot200Error(error)) {
    switch (error.response.status) {
      case 404:
        return (
          <Container fluid className="pt-[10%]">
            <Error404 />
          </Container>
        )
      case 401:
        logout()
        navigate({ to: "/" })
        break
    }
  }

  const handleReturn = () => {
    reset()
    navigate({ to: "/" })
  }

  return (
    <div className="flex flex-col justify-center items-center h-full pt-[10%]">
      <p className="font-semibold text-primary">
        500
      </p>
      <Title order={1} className="mt-4 text-balance font-semibold tracking-tight">
        Server Error
      </Title>
      <p className="mt-6 text-pretty text-lg font-medium text-gray-500 sm:text-xl/8">
        <span>Something went wrong</span>
      </p>
      <div className="mt-10">
        <Button onClick={handleReturn}>
          Go back
        </Button>
      </div>
    </div>
  )
}
