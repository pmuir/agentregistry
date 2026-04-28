"use client"

import { useState } from "react"
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { createAgentGateway } from "@/lib/admin-api"

interface AddAgentGatewayDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onGatewayAdded: () => void
}

export function AddAgentGatewayDialog({ open, onOpenChange, onGatewayAdded }: AddAgentGatewayDialogProps) {
  const [id, setId] = useState("")
  const [name, setName] = useState("")
  const [address, setAddress] = useState("")
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)

    try {
      if (!id.trim()) {
        throw new Error("Gateway ID is required")
      }
      if (!name.trim()) {
        throw new Error("Gateway name is required")
      }
      if (!address.trim()) {
        throw new Error("Gateway address is required")
      }

      try {
        new URL(address.trim())
      } catch {
        throw new Error("Address must be a valid URL (e.g., http://localhost:8081)")
      }

      await createAgentGateway({
        body: {
          id: id.trim(),
          name: name.trim(),
          address: address.trim(),
        },
        throwOnError: true,
      })

      setId("")
      setName("")
      setAddress("")

      onGatewayAdded()
      onOpenChange(false)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to add gateway")
    } finally {
      setLoading(false)
    }
  }

  const handleCancel = () => {
    setId("")
    setName("")
    setAddress("")
    setError(null)
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>Add Agent Gateway</DialogTitle>
          <DialogDescription>
            Register a running agentgateway instance
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4 py-4">
          <div className="space-y-2">
            <Label htmlFor="gw-id">
              ID <span className="text-red-500">*</span>
            </Label>
            <Input
              id="gw-id"
              placeholder="my-gateway"
              value={id}
              onChange={(e) => setId(e.target.value)}
              disabled={loading}
              required
            />
            <p className="text-xs text-muted-foreground">
              A unique identifier for this gateway instance
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="gw-name">
              Name <span className="text-red-500">*</span>
            </Label>
            <Input
              id="gw-name"
              placeholder="Production Gateway"
              value={name}
              onChange={(e) => setName(e.target.value)}
              disabled={loading}
              required
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="gw-address">
              Address <span className="text-red-500">*</span>
            </Label>
            <Input
              id="gw-address"
              placeholder="http://localhost:8081"
              value={address}
              onChange={(e) => setAddress(e.target.value)}
              disabled={loading}
              required
            />
            <p className="text-xs text-muted-foreground">
              The URL of the running agentgateway (health checked automatically)
            </p>
          </div>

          {error && (
            <div className="rounded-md bg-red-50 p-3 text-sm text-red-800">
              {error}
            </div>
          )}

          <div className="flex justify-end gap-2">
            <Button type="button" variant="outline" onClick={handleCancel} disabled={loading}>
              Cancel
            </Button>
            <Button type="submit" disabled={loading}>
              {loading ? "Adding..." : "Add Gateway"}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
