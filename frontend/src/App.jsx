import { useState, useEffect, useRef } from 'react'
import './App.css'

const CLIENT_ID = '11111111-1111-1111-1111-111111111111'
const ADMIN_ID = '22222222-2222-2222-2222-222222222222'

const statuses = [
  'PENDING',
  'ASSIGNED',
  'ON_THE_WAY',
  'IN_PROGRESS',
  'COMPLETED',
]

function App() {
  const [orderId, setOrderId] = useState('')
  const [order, setOrder] = useState(null)
  const [events, setEvents] = useState([])
  const [serviceName, setServiceName] = useState('')
  const [description, setDescription] = useState('')
  const [loading, setLoading] = useState(false)
  const [message, setMessage] = useState('')
  const [connected, setConnected] = useState(false)

  const socketRef = useRef(null)

  // Ambil detail pesanan dan riwayat
  const loadOrder = async (id) => {
    try {
      const response = await fetch(`/api/orders/${id}`)
      const result = await response.json()

      if (!response.ok) {
        throw new Error(result.error || 'Pesanan tidak ditemukan')
      }

      setOrder(result.data)

      const eventResponse = await fetch(`/api/orders/${id}/events`)
      const eventResult = await eventResponse.json()

      if (eventResponse.ok) {
        setEvents(eventResult.data || [])
      }

      setMessage('')
    } catch (error) {
      setMessage(error.message)
      setOrder(null)
      setEvents([])
    }
  }

  // Buat pesanan baru
  const createOrder = async (e) => {
    e.preventDefault()

    setLoading(true)
    setMessage('')

    try {
      const response = await fetch('/api/orders', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          client_id: CLIENT_ID,
          service_name: serviceName,
          description: description,
        }),
      })

      const result = await response.json()

      if (!response.ok) {
        throw new Error(result.error || 'Gagal membuat pesanan')
      }

      setOrderId(result.data.ID)
      await loadOrder(result.data.ID)

      setServiceName('')
      setDescription('')

      setMessage('Pesanan berhasil dibuat!')
    } catch (error) {
      setMessage(error.message)
    } finally {
      setLoading(false)
    }
  }

  // Koneksi WebSocket
  useEffect(() => {
    if (!order?.ID) return

    const socket = new WebSocket(
      `ws://${window.location.host}/ws/orders/${order.ID}`
    )

    socketRef.current = socket

    socket.onopen = () => {
      setConnected(true)
    }

    socket.onmessage = (event) => {
      const newEvent = JSON.parse(event.data)

      setEvents((previous) => {
        const exists = previous.some(
          (item) => item.Sequence === newEvent.Sequence
        )

        if (exists) return previous

        return [...previous, newEvent]
      })

      setOrder((previous) => {
        if (!previous) return previous

        return {
          ...previous,
          Status: newEvent.NewStatus,
          Version: previous.Version + 1,
          LastSequence: newEvent.Sequence,
        }
      })
    }

    socket.onclose = () => {
      setConnected(false)
    }

    socket.onerror = () => {
      setConnected(false)
    }

    return () => {
      socket.close()
      setConnected(false)
    }
  }, [order?.ID])

  // Cari pesanan berdasarkan ID
  const searchOrder = async (e) => {
    e.preventDefault()

    if (!orderId.trim()) return

    await loadOrder(orderId.trim())
  }

  // Simulasi perubahan status oleh Admin
  const updateStatus = async () => {
    if (!order) return

    const currentIndex = statuses.indexOf(order.Status)

    if (currentIndex === -1 || currentIndex >= statuses.length - 1) {
      return
    }

    const nextStatus = statuses[currentIndex + 1]

    try {
      const response = await fetch(`/api/orders/${order.ID}/status`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          new_status: nextStatus,
          actor_id: ADMIN_ID,
          expected_version: order.Version,
        }),
      })

      const result = await response.json()

      if (!response.ok) {
        throw new Error(result.error || 'Gagal memperbarui status')
      }

      setMessage('Status berhasil diperbarui!')
      
      // Ambil data terbaru dari backend
      await loadOrder(order.ID)
    } catch (error) {
      setMessage(error.message)
    }
  }

  const formatStatus = (status) => {
    return status?.replaceAll('_', ' ')
  }

  return (
    <div className="app">
      <header className="header">
        <div>
          <h1>OrderTrack</h1>
          <p>Real-Time Order & Job Tracking System</p>
        </div>

        <div className={`connection ${connected ? 'online' : 'offline'}`}>
          <span></span>
          {connected ? 'Live Connected' : 'Disconnected'}
        </div>
      </header>

      <main className="container">

        <section className="hero">
          <h2>Track Your Service Order</h2>
          <p>Pantau status pesanan layanan secara real-time.</p>
        </section>

        <div className="grid">

          <section className="card">
            <h3>Buat Pesanan Baru</h3>

            <form onSubmit={createOrder}>
              <label>Nama Layanan</label>

              <input
                type="text"
                placeholder="Contoh: Perbaikan AC"
                value={serviceName}
                onChange={(e) => setServiceName(e.target.value)}
                required
              />

              <label>Deskripsi Kerusakan</label>

              <textarea
                placeholder="Jelaskan masalah yang dialami..."
                value={description}
                onChange={(e) => setDescription(e.target.value)}
              />

              <button type="submit" disabled={loading}>
                {loading ? 'Memproses...' : 'Buat Pesanan'}
              </button>
            </form>
          </section>

          <section className="card">
            <h3>Cari Pesanan</h3>

            <form onSubmit={searchOrder}>
              <label>Order ID</label>

              <input
                type="text"
                placeholder="Masukkan ID pesanan"
                value={orderId}
                onChange={(e) => setOrderId(e.target.value)}
              />

              <button type="submit">
                Cari Pesanan
              </button>
            </form>

            {message && (
              <div className="message">
                {message}
              </div>
            )}
          </section>

        </div>

        {order && (
          <>
            <section className="card order-card">
              <div className="order-heading">
                <div>
                  <h3>Detail Pesanan</h3>
                  <p className="order-id">{order.ID}</p>
                </div>

                <span className="status-badge">
                  {formatStatus(order.Status)}
                </span>
              </div>

              <div className="details">
                <div>
                  <small>Nama Layanan</small>
                  <strong>{order.ServiceName}</strong>
                </div>

                <div>
                  <small>Deskripsi</small>
                  <strong>{order.Description || '-'}</strong>
                </div>

                <div>
                  <small>Version</small>
                  <strong>{order.Version}</strong>
                </div>

                <div>
                  <small>Last Sequence</small>
                  <strong>{order.LastSequence}</strong>
                </div>
              </div>
            </section>

            <section className="card">
              <h3>Tracking Progress</h3>

              <div className="timeline">
                {statuses.map((status, index) => {
                  const currentIndex = statuses.indexOf(order.Status)

                  const completed = index <= currentIndex

                  return (
                    <div
                      className={`timeline-item ${completed ? 'completed' : ''}`}
                      key={status}
                    >
                      <div className="timeline-dot">
                        {completed ? '✓' : index + 1}
                      </div>

                      <div>
                        <strong>{formatStatus(status)}</strong>
                        <small>
                          {index === currentIndex
                            ? 'Status saat ini'
                            : completed
                              ? 'Selesai'
                              : 'Menunggu'}
                        </small>
                      </div>
                    </div>
                  )
                })}
              </div>

              {order.Status !== 'COMPLETED' &&
                order.Status !== 'CANCELLED' && (
                  <button onClick={updateStatus}>
                    Simulasi Update Status →
                  </button>
                )}
            </section>

            <section className="card">
              <h3>Riwayat Aktivitas</h3>

              {events.length === 0 ? (
                <p>Belum ada riwayat aktivitas.</p>
              ) : (
                <div className="history">
                  {events.map((event) => (
                    <div className="history-item" key={event.ID}>
                      <div className="history-dot"></div>

                      <div>
                        <strong>
                          {formatStatus(event.NewStatus)}
                        </strong>

                        <p>{event.Note}</p>

                        <small>
                          Sequence #{event.Sequence} •{' '}
                          {new Date(event.CreatedAt).toLocaleString('id-ID')}
                        </small>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </section>
          </>
        )}

      </main>

      <footer>
        OrderTrack © 2026 — Real-Time Service Management
      </footer>
    </div>
  )
}

export default App