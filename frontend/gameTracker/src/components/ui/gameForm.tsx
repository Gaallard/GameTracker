import { useEffect, useState } from "react"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Textarea } from "@/components/ui/textarea"
import { createGame, updateGame, type Game } from "@/services/api"

interface GameFormProps {
    onGameAdded: (game: Game) => void
    gameToEdit?: Game | null
    onEditComplete?: () => void
}

export default function GameForm({ onGameAdded, gameToEdit, onEditComplete }: GameFormProps) {
    const [title, setTitle] = useState("")
    const [platform, setPlatform] = useState("")
    const [genre, setGenre] = useState("")
    const [status, setStatus] = useState("Backlog")
    const [note, setNote] = useState("")

    useEffect(() => {
        if (gameToEdit) {
            setTitle(gameToEdit.title)
            setPlatform(gameToEdit.platform)
            setGenre(gameToEdit.genre)
            setStatus(gameToEdit.status)
            setNote(gameToEdit.personalNote)
        } else {
            setTitle("")
            setPlatform("")
            setGenre("")
            setStatus("Backlog")
            setNote("")
        }
    }, [gameToEdit])

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        const gameData = {
            title,
            platform,
            genre,
            status,
            progress: gameToEdit ? gameToEdit.progress : 0,
            hoursPlayed: gameToEdit ? gameToEdit.hoursPlayed : 0,
            personalNote: note,
            score: gameToEdit ? gameToEdit.score : 0,
            startedAt: gameToEdit ? gameToEdit.startedAt : new Date().toISOString(),
            finishedAt: gameToEdit ? gameToEdit.finishedAt : new Date().toISOString(),
            coverURL: gameToEdit ? gameToEdit.coverURL : "",
        }

        try {
            let res
            if (gameToEdit) {
                res = await updateGame(gameToEdit.id, gameData)
                alert("Juego actualizado con éxito")
            } else {
                res = await createGame(gameData)
                alert("Juego creado con éxito")
            }
            onGameAdded(res.data)
            if (onEditComplete) onEditComplete()

            if (!gameToEdit) {
                // Limpiar formulario sólo si estamos creando
                setTitle("")
                setPlatform("")
                setGenre("")
                setStatus("Backlog")
                setNote("")
            }
        } catch (err) {
            alert("Error al guardar juego")
            console.error(err)
        }
    }

    return (
        <form onSubmit={handleSubmit} className="max-w-md mx-auto space-y-4 p-4 border rounded-xl shadow bg-white">
            <h2 className="text-xl font-bold text-center">
                {gameToEdit ? "Editar Juego" : "Agregar Juego"}
            </h2>

            <Input placeholder="Título" value={title} onChange={e => setTitle(e.target.value)} required />
            <Input placeholder="Plataforma" value={platform} onChange={e => setPlatform(e.target.value)} required />
            <Input placeholder="Género" value={genre} onChange={e => setGenre(e.target.value)} required />

            <div>
                <label className="block mb-1 text-sm font-medium text-gray-700">Estado</label>
                <select
                    value={status}
                    onChange={e => setStatus(e.target.value)}
                    className="w-full p-2 border rounded text-sm"
                >
                    <option value="Backlog">Backlog</option>
                    <option value="Playing">Playing</option>
                    <option value="Completed">Completed</option>
                    <option value="Dropped">Dropped</option>
                </select>
            </div>

            <Textarea placeholder="Nota personal" value={note} onChange={e => setNote(e.target.value)} />

            <Button type="submit" className="w-full">
                {gameToEdit ? "Actualizar juego" : "Guardar"}
            </Button>
        </form>
    )
}
