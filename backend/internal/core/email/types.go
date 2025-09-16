package email

// EmailPayload follows Maileroo's API format.
type EmailPayload struct {
	FromEmail   EmailObject   `json:"from"`
	ToEmail     []EmailObject `json:"to"`
	Subject     string        `json:"subject"`
	HtmlBody    string        `json:"html"`
	Attachments []File        `json:"attachments"`
}

type EmailObject struct {
	Address     string `json:"address"`
	DisplayName string `json:"display_name"`
}

type File struct {
	Name        string `json:"file_name"`
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
}

const EmailHTMLTemplate = `
<!DOCTYPE html>
<html lang="es">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Nueva Compra</title>
    <style>
      body {
        font-family: Arial, sans-serif;
        background-color: #f6f6f6;
        color: #333333;
        padding: 20px;
        margin: 0;
      }
      .container {
        background-color: #ffffff;
        padding: 20px;
        max-width: 600px;
        margin: auto;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      }
      h1 {
        color: #e67e22;
        font-size: 20px;
      }
      .details {
        margin-top: 20px;
      }
      .details p {
        margin: 8px 0;
      }
      .footer {
        margin-top: 30px;
        font-size: 12px;
        color: #999;
        text-align: center;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <h1>📩 Nueva Compra Realizada</h1>
      <p>
        Se ha recibido una nueva solicitud de compra que requiere verificación:
      </p>

      <div class="details">
        <p><strong>👤 Usuario ID:</strong> %s</p>
        <p><strong>📅 Fecha:</strong> %s</p>
        <p><strong>🎟️ Cantidad de boletos:</strong> %d</p>
        <p><strong>💵 Monto en USD:</strong> $%.2f</p>
        <p><strong>💴 Monto en Bs:</strong> %.2f Bs</p>
        <p><strong>💳 Método de pago:</strong> %s</p>
        <p><strong>🔢 Últimos dígitos:</strong> %s</p>
      </div>

      <p style="margin-top: 20px">
        🖼️ Se adjuntó la captura del pago como imagen.
      </p>

      <div class="footer">
        Este mensaje fue generado automáticamente por el sistema de rifas.
      </div>
    </div>
  </body>
</html>
`
