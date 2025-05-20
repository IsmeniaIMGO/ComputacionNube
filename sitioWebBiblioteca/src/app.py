from flask import Flask, render_template, request, redirect, url_for, make_response
import pymysql

app = Flask(__name__)

# Configuración de conexión a MySQL
DB_CONFIG = {
    'host': 'localhost',
    'user': 'root',
    'password': 'root',
    'database': 'biblioteca',
    'charset': 'utf8mb4',
    'cursorclass': pymysql.cursors.DictCursor
}

def get_connection():
    return pymysql.connect(**DB_CONFIG)

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/libros')
def libros():
    conn = get_connection()
    with conn:
        with conn.cursor() as cur:
            cur.execute("SELECT l.id, l.titulo, a.nombre AS autor, c.nombre AS categoria FROM libro l "
                        "JOIN autor a ON l.autor_id = a.id "
                        "JOIN categoria c ON l.categoria_id = c.id")
            libros = cur.fetchall()
    return render_template('libros.html', libros=libros)

@app.route('/autores')
def autores():
    conn = get_connection()
    with conn:
        with conn.cursor() as cur:
            cur.execute("SELECT id, nombre FROM autor")
            autores = cur.fetchall()
    return render_template('autores.html', autores=autores)

@app.route('/categorias')
def categorias():
    conn = get_connection()
    with conn:
        with conn.cursor() as cur:
            cur.execute("SELECT id, nombre FROM categoria")
            categorias = cur.fetchall()
    return render_template('categorias.html', categorias=categorias)

@app.route('/agregar_libro', methods=['GET', 'POST'])
def agregar_libro():
    conn = get_connection()
    with conn:
        with conn.cursor() as cur:
            cur.execute("SELECT id, nombre FROM autor")
            autores = cur.fetchall()
            cur.execute("SELECT id, nombre FROM categoria")
            categorias = cur.fetchall()

        if request.method == 'POST':
            titulo = request.form['titulo']
            autor_id = request.form['autor_id']
            categoria_id = request.form['categoria_id']

            if titulo and autor_id and categoria_id:
                with conn.cursor() as cur:
                    cur.execute(
                        "INSERT INTO libro (titulo, autor_id, categoria_id) VALUES (%s, %s, %s)",
                        (titulo, autor_id, categoria_id)
                    )
                conn.commit()
                return redirect(url_for('libros'))

    return render_template('agregar_libro.html', autores=autores, categorias=categorias)


@app.route('/agregar_autor', methods=['GET', 'POST'])
def agregar_autor():
    if request.method == 'POST':
        nombre = request.form.get('nombre')  # Usa .get() para evitar errores si falta el campo

        if nombre:  # Verifica que se haya enviado un nombre válido
            conn = get_connection()
            with conn:
                with conn.cursor() as cur:
                    cur.execute("INSERT INTO autor (nombre) VALUES (%s)", (nombre,))
                conn.commit()
            return redirect(url_for('autores'))
        else:
            # Aquí podrías agregar un mensaje de error si quieres validación
            pass
    return render_template('agregar_autor.html')

@app.route('/agregar_categoria', methods=['GET', 'POST'])
def agregar_categoria():
    if request.method == 'POST':
        nombre = request.form.get('nombre', '').strip()

        # Validación básica del campo
        if nombre:
            conn = get_connection()
            with conn:
                with conn.cursor() as cur:
                    cur.execute("INSERT INTO categoria (nombre) VALUES (%s)", (nombre,))
                conn.commit()
            return redirect(url_for('categorias'))
        else:
            # Opcional: manejar el error si nombre está vacío
            return render_template('agregar_categoria.html', error="El nombre no puede estar vacío.")
    
    return render_template('agregar_categoria.html')


@app.route('/editar_libro/<int:libro_id>', methods=['GET', 'POST'])
def editar_libro(libro_id):
    conn = get_connection()
    with conn:
        with conn.cursor() as cur:
            # Obtener el libro por id
            cur.execute("SELECT id, titulo, autor_id, categoria_id FROM libro WHERE id = %s", (libro_id,))
            libro = cur.fetchone()

            # Obtener autores y categorías para el select
            cur.execute("SELECT id, nombre FROM autor")
            autores = cur.fetchall()
            cur.execute("SELECT id, nombre FROM categoria")
            categorias = cur.fetchall()

        if request.method == 'POST':
            titulo = request.form.get('titulo', '').strip()
            autor_id = request.form.get('autor_id', '').strip()
            categoria_id = request.form.get('categoria_id', '').strip()

            # Validar que los campos no estén vacíos
            if titulo and autor_id and categoria_id:
                with conn.cursor() as cur:
                    cur.execute(
                        "UPDATE libro SET titulo = %s, autor_id = %s, categoria_id = %s WHERE id = %s",
                        (titulo, autor_id, categoria_id, libro_id)
                    )
                conn.commit()
                return redirect(url_for('libros'))
            else:
                # Podrías manejar un error si algún campo está vacío
                error = "Todos los campos son obligatorios."
                return render_template('editar_libro.html', libro=libro, autores=autores, categorias=categorias, error=error)

    # GET request o si falla validación, mostrar formulario con datos actuales
    return render_template('editar_libro.html', libro=libro, autores=autores, categorias=categorias)


@app.route('/editar_autor/<int:autor_id>', methods=['GET', 'POST'])
def editar_autor(autor_id):
    conn = get_connection()
    error = None
    with conn:
        with conn.cursor() as cur:
            # Obtener el autor por id
            cur.execute("SELECT id, nombre FROM autor WHERE id = %s", (autor_id,))
            autor = cur.fetchone()

            if request.method == 'POST':
                nombre = request.form.get('nombre', '').strip()

                if nombre:
                    cur.execute(
                        "UPDATE autor SET nombre = %s WHERE id = %s",
                        (nombre, autor_id)
                    )
                    conn.commit()
                    return redirect(url_for('autores'))
                else:
                    error = "El nombre no puede estar vacío."

    return render_template('editar_autor.html', autor=autor, error=error)

@app.route('/editar_categoria/<int:categoria_id>', methods=['GET', 'POST'])
def editar_categoria(categoria_id):
    conn = get_connection()
    error = None
    with conn:
        with conn.cursor() as cur:
            # Obtener la categoría por id
            cur.execute("SELECT id, nombre FROM categoria WHERE id = %s", (categoria_id,))
            categoria = cur.fetchone()

        if request.method == 'POST':
            nombre = request.form.get('nombre', '').strip()

            if nombre:
                with conn.cursor() as cur:
                    cur.execute(
                        "UPDATE categoria SET nombre = %s WHERE id = %s",
                        (nombre, categoria_id)
                    )
                conn.commit()
                return redirect(url_for('categorias'))
            else:
                error = "El nombre no puede estar vacío."

    return render_template('editar_categoria.html', categoria=categoria, error=error)


@app.route('/eliminar_libro/<int:libro_id>', methods=['POST'])
def eliminar_libro(libro_id):
    conn = get_connection()
    with conn:
        with conn.cursor() as cur:
            cur.execute("DELETE FROM libro WHERE id = %s", (libro_id,))
        conn.commit()
    return redirect(url_for('libros'))


@app.route('/eliminar_autor/<int:autor_id>', methods=['GET', 'POST'])
def eliminar_autor(autor_id):
    conn = get_connection()
    error = None
    with conn:
        with conn.cursor() as cur:
            # Consultar autor
            cur.execute("SELECT id, nombre FROM autor WHERE id = %s", (autor_id,))
            autor = cur.fetchone()

            if request.method == 'POST':
                # Verificar si hay libros asociados
                cur.execute("SELECT COUNT(*) AS total FROM libro WHERE autor_id = %s", (autor_id,))
                resultado = cur.fetchone()
                if resultado['total'] > 0:
                    error = "No se puede eliminar el autor porque tiene libros asociados."
                else:
                    cur.execute("DELETE FROM autor WHERE id = %s", (autor_id,))
                    conn.commit()
                    return redirect(url_for('autores'))

    return render_template('eliminar_autor.html', autor=autor, error=error)

@app.route('/eliminar_categoria/<int:categoria_id>', methods=['GET', 'POST'])
def eliminar_categoria(categoria_id):
    conn = get_connection()
    error = None

    with conn:
        with conn.cursor() as cur:
            # Obtener la categoría a eliminar
            cur.execute("SELECT id, nombre FROM categoria WHERE id = %s", (categoria_id,))
            categoria = cur.fetchone()

            if request.method == 'POST':
                # Verificar si la categoría tiene libros asociados
                cur.execute("SELECT COUNT(*) AS total FROM libro WHERE categoria_id = %s", (categoria_id,))
                resultado = cur.fetchone()
                
                if resultado['total'] > 0:
                    error = "No se puede eliminar la categoría porque tiene libros asociados."
                else:
                    # No hay libros asociados, eliminar categoría
                    cur.execute("DELETE FROM categoria WHERE id = %s", (categoria_id,))
                    conn.commit()
                    return redirect(url_for('categorias'))

    return render_template('eliminar_categoria.html', categoria=categoria, error=error)

if __name__ == '__main__':
    app.run(debug=True)
