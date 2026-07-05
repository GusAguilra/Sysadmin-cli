# Recursos visuales de sysadmin-cli

Archivos de promoción e imágenes para el proyecto.

## Archivos

### `banner.svg`
Banner vectorial para GitHub, blog, o redes sociales. Contiene:
- Logo y nombre del proyecto
- Descripción breve
- Características principales
- Información de licencia

**Cómo usar:**
- Abrir en navegador web
- Convertir a PNG: `inkscape banner.svg -o banner.png`
- Exportar a PDF: `inkscape banner.svg --export-pdf=banner.pdf`

### `demo.txt`
Demostración interactiva del TUI mostrando:
- Vista inicial (lista de categorías)
- Vista de categoría (comandos)
- Vista detalle (comando completo)
- Modo búsqueda
- Uso en línea de comandos

**Cómo usar:**
- Copiar en documentación
- Mostrar en README.md
- Usar como referencia para screenshots reales

## Próximos pasos

Para mejorar estos recursos:

1. **GIF animado**: Usar `asciinema` para grabar sesión interactiva
   ```bash
   asciinema rec demo.cast
   agg demo.cast demo.gif
   ```

2. **PNG del banner**: Convertir SVG a PNG
   ```bash
   inkscape banner.svg -o banner.png
   ```

3. **Screenshots reales**: Capturar pantallazos de la interfaz
   ```bash
   import -window root screenshot.png
   ```

## Notas

- Todos los recursos están diseñados para uso en GitHub (README, releases, etc.)
- El SVG es escalable y se verá bien en cualquier resolución
- El demo.txt proporciona una vista textual del flujo de usuario
